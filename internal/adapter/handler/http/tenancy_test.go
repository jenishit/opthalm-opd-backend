package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

// TestTenancy_CrossClinicIsolation is the load-bearing test of the whole
// multi-tenancy feature: clinic A's data must be completely invisible to
// clinic B, across every module — not just absent from lists, but a 404
// (never a data leak) when B tries to fetch A's records directly by ID.
func TestTenancy_CrossClinicIsolation(t *testing.T) {
	ts := testutil.NewTestServer(t)
	tokenA, _ := ts.AdminToken(t)
	tokenB, _ := ts.SecondClinicAdminToken(t)

	// Clinic A creates a patient, an inventory item, and (via that item) an invoice.
	patientA := ts.CreatePatient(t, tokenA, "Clinic A Patient", "9811111111")

	var itemA struct {
		ID uuid.UUID `json:"id"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/inventory/items", tokenA, map[string]any{
		"category": "frame", "sku": "TENANCY-A-1", "name": "A's Frame",
		"cost_price": 1000, "selling_price": 2000, "quantity_on_hand": 10, "reorder_threshold": 3,
	}, &itemA)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var invoiceA struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/billing/invoice", tokenA, map[string]any{
		"patient_id": patientA,
		"items": []map[string]any{
			{"item_type": "frame", "description": "A's Frame", "inventory_item_id": itemA.ID, "quantity": 1, "unit_price": 2000},
		},
	}, &invoiceA)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var vendorA struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/inventory/vendors", tokenA, map[string]any{"name": "A's Vendor"}, &vendorA)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var labJobA struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/lab-jobs", tokenA, map[string]any{
		"patient_id": patientA, "job_type": "lens_fitting",
	}, &labJobA)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Clinic B must not be able to fetch any of clinic A's records by ID.
	notFoundChecks := map[string]string{
		"patient":        "/api/patient/" + patientA.String(),
		"inventory item": "/api/inventory/items/" + itemA.ID.String(),
		"invoice":        "/api/billing/invoice/" + invoiceA.ID.String(),
		"vendor":         "/api/inventory/vendors/" + vendorA.ID.String(),
		"lab job":        "/api/lab-jobs/" + labJobA.ID.String(),
	}
	for label, path := range notFoundChecks {
		resp := ts.Do(t, http.MethodGet, path, tokenB, nil, nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "clinic B fetching clinic A's %s by ID must 404, not leak data", label)
	}

	// Clinic B's list/search endpoints must not include clinic A's records either.
	var patients []struct{ ID uuid.UUID }
	resp = ts.DoData(t, http.MethodGet, "/api/patient", tokenB, nil, &patients)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, patients, "clinic B's patient list must not include clinic A's patient")

	var items []struct{ ID uuid.UUID }
	resp = ts.DoData(t, http.MethodGet, "/api/inventory/items", tokenB, nil, &items)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, items, "clinic B's inventory list must not include clinic A's item")

	var invoices []struct{ ID uuid.UUID }
	resp = ts.DoData(t, http.MethodGet, "/api/billing/invoice", tokenB, nil, &invoices)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, invoices, "clinic B's invoice list must not include clinic A's invoice")

	// Clinic B can create its own patient with the exact same phone number —
	// proves there's no accidental global-uniqueness leak between tenants.
	patientB := ts.CreatePatient(t, tokenB, "Clinic B Patient", "9811111111")
	assert.NotEqual(t, patientA, patientB)

	// Clinic B cannot record a payment against clinic A's invoice by guessing its ID.
	resp = ts.Do(t, http.MethodPost, "/api/billing/invoice/"+invoiceA.ID.String()+"/payments", tokenB, map[string]any{
		"amount": 100, "method": "cash",
	}, nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "clinic B must not be able to pay against clinic A's invoice")
}

func TestTenancy_UsersAreScopedToTheirOwnClinic(t *testing.T) {
	ts := testutil.NewTestServer(t)
	tokenA, _ := ts.AdminToken(t)
	ts.EnsureRole(t, "ROLE_BILLING")

	// Clinic A creates a ROLE_BILLING staff user via the real (admin-gated) endpoint.
	resp := ts.Do(t, http.MethodPost, "/api/user/create", tokenA, map[string]string{
		"first_name": "Billing", "last_name": "Clerk",
		"email": "clerk@clinic-a.local", "password": testutil.TestPassword, "role_name": "ROLE_BILLING",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var login struct {
		AccessToken string `json:"access_token"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email": "clerk@clinic-a.local", "password": testutil.TestPassword,
	}, &login)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// That new user's token should carry clinic A's ID, so a patient they
	// create lands in clinic A, not some default/empty tenant.
	patientID := ts.CreatePatient(t, login.AccessToken, "Clerk's Patient", "9822222222")

	var fetched struct {
		PatientID uuid.UUID `json:"patient_id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/patient/"+patientID.String(), tokenA, nil, &fetched)
	require.Equal(t, http.StatusOK, resp.StatusCode, "clinic A's admin should see the patient the new clerk created")
	assert.Equal(t, patientID, fetched.PatientID)
}
