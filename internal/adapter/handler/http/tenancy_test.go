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
