package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
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
