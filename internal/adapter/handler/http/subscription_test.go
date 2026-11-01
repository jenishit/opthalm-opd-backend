package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
)

func TestSubscription_CancelledClinicGets402OnTenantRoutesButCanStillLogin(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	ts.SetSubscriptionStatus(t, ts.ClinicID, "cancelled", time.Now().Add(24*time.Hour))

	resp := ts.Do(t, http.MethodGet, "/api/patient", token, nil, nil)
	assert.Equal(t, http.StatusPaymentRequired, resp.StatusCode)

	// Login must still work — the client needs to be able to authenticate to
	// even find out the subscription is the problem.
	var login struct {
		AccessToken string `json:"access_token"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email":    "admin@test.local",
		"password": testutil.TestPassword,
	}, &login)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, login.AccessToken)
}

func TestSubscription_ExpiredPeriodEndGets402(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	// Status still "active" but the period end has already passed.
	ts.SetSubscriptionStatus(t, ts.ClinicID, "active", time.Now().Add(-24*time.Hour))

	resp := ts.Do(t, http.MethodGet, "/api/billing/invoice", token, nil, nil)
	assert.Equal(t, http.StatusPaymentRequired, resp.StatusCode)
}

func TestSubscription_TrialingClinicIsAllowed(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	// NewTestServer already seeds a trialing subscription — no setup needed.
	resp := ts.Do(t, http.MethodGet, "/api/patient", token, nil, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSubscription_CalculatorsWorkRegardlessOfSubscriptionStatus(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	ts.SetSubscriptionStatus(t, ts.ClinicID, "cancelled", time.Now().Add(24*time.Hour))

	resp := ts.Do(t, http.MethodPost, "/api/calculators/spherical-equivalent", token, map[string]any{
		"sphere": 1.0, "cylinder": 1.0,
	}, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "free tools should work even for a cancelled subscription")
}

func TestSubscription_PlatformOperatorCanReactivate(t *testing.T) {
	ts := testutil.NewTestServer(t)
	clinicToken, _ := ts.AdminToken(t)
	superToken, _ := ts.SuperadminToken(t)

	ts.SetSubscriptionStatus(t, ts.ClinicID, "cancelled", time.Now().Add(24*time.Hour))

	resp := ts.Do(t, http.MethodGet, "/api/patient", clinicToken, nil, nil)
	require.Equal(t, http.StatusPaymentRequired, resp.StatusCode)

	// A clinic's own admin cannot reactivate itself.
	resp = ts.Do(t, http.MethodPut, "/api/platform/subscriptions/"+ts.ClinicID.String(), clinicToken, map[string]any{
		"plan_name": "pro", "status": "active", "current_period_end": time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "a clinic admin must not be able to reactivate its own subscription")

	// Only the platform superadmin can.
	resp = ts.Do(t, http.MethodPut, "/api/platform/subscriptions/"+ts.ClinicID.String(), superToken, map[string]any{
		"plan_name": "pro", "status": "active", "current_period_end": time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodGet, "/api/patient", clinicToken, nil, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "clinic should regain access after the platform operator reactivates it")
}

func TestSubscription_ClinicBillingViewShowsLiveUsageAgainstPlatformSetLimits(t *testing.T) {
	ts := testutil.NewTestServer(t)
	clinicToken, _ := ts.AdminToken(t)
	superToken, _ := ts.SuperadminToken(t)

	// The platform operator sets the plan's limits...
	resp := ts.Do(t, http.MethodPut, "/api/platform/subscriptions/"+ts.ClinicID.String(), superToken, map[string]any{
		"plan_name":             "pro",
		"status":                "active",
		"current_period_end":    time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
		"seats_limit":           5,
		"patient_records_limit": 100,
		"locations_limit":       2,
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// A clinic admin can never reach the platform-only endpoint...
	resp = ts.Do(t, http.MethodGet, "/api/platform/subscriptions/"+ts.ClinicID.String(), clinicToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// ...but can read its own billing view, which carries those limits plus
	// live usage computed from this clinic's actual rows.
	var billing struct {
		PlanName            string  `json:"plan_name"`
		SeatsLimit          int     `json:"seats_limit"`
		SeatsUsed           int     `json:"seats_used"`
		PatientRecordsLimit int     `json:"patient_records_limit"`
		PatientRecordsUsed  int     `json:"patient_records_used"`
		LocationsLimit      int     `json:"locations_limit"`
		LocationsUsed       int     `json:"locations_used"`
		PaymentMethodLabel  *string `json:"payment_method_label"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/billing/subscription", clinicToken, nil, &billing)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pro", billing.PlanName)
	assert.Equal(t, 5, billing.SeatsLimit)
	assert.Equal(t, 100, billing.PatientRecordsLimit)
	assert.Equal(t, 2, billing.LocationsLimit)
	seatsBefore := billing.SeatsUsed
	assert.GreaterOrEqual(t, seatsBefore, 1) // the admin user created by AdminToken
	assert.Equal(t, 0, billing.PatientRecordsUsed)
	assert.Equal(t, 0, billing.LocationsUsed)
	assert.Nil(t, billing.PaymentMethodLabel) // nothing syncs this yet

	// Usage reflects real data, live — no caching, no drift.
	ts.CreatePatient(t, clinicToken, "Billing Usage Patient", "9800000040")
	resp = ts.Do(t, http.MethodPost, "/api/admin/clinic/locations", clinicToken, map[string]any{
		"name": "Second Room", "exam_rooms": 1,
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/billing/subscription", clinicToken, nil, &billing)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, billing.PatientRecordsUsed)
	assert.Equal(t, 1, billing.LocationsUsed)
	assert.Equal(t, seatsBefore, billing.SeatsUsed) // unchanged, no new users were created
}
