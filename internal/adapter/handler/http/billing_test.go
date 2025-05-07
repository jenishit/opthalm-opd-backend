package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func createInventoryItem(t *testing.T, ts *testutil.TestServer, token, sku string, qty int) uuid.UUID {
	t.Helper()
	var out struct {
		ID uuid.UUID `json:"id"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/inventory/items", token, map[string]any{
		"category":          "frame",
		"sku":               sku,
		"name":              "Test Frame",
		"cost_price":        1000,
		"selling_price":     2000,
		"quantity_on_hand":  qty,
		"reorder_threshold": 3,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return out.ID
}

func TestBilling_CreateInvoiceDeductsStock(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Billing Patient", "9800000020")
	itemID := createInventoryItem(t, ts, token, "FRM-BILL-1", 10)

	var invoice struct {
		ID          uuid.UUID `json:"id"`
		InvoiceNo   string    `json:"invoice_no"`
		TotalAmount float64   `json:"total_amount"`
		DueAmount   float64   `json:"due_amount"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/billing/invoice", token, map[string]any{
		"patient_id": patientID,
		"items": []map[string]any{
			{"item_type": "frame", "description": "Test Frame", "inventory_item_id": itemID, "quantity": 2, "unit_price": 2000},
		},
	}, &invoice)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, invoice.InvoiceNo)
	assert.Equal(t, float64(4000), invoice.TotalAmount)
	assert.Equal(t, float64(4000), invoice.DueAmount)

	var item struct {
		QuantityOnHand int `json:"quantity_on_hand"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/inventory/items/"+itemID.String(), token, nil, &item)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 8, item.QuantityOnHand, "stock should be decremented by the invoice quantity")
}

func TestBilling_PartialPaymentThenOverpayRejected(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Billing Patient 2", "9800000021")

	var invoice struct {
		ID uuid.UUID `json:"id"`
	}
	ts.DoData(t, http.MethodPost, "/api/billing/invoice", token, map[string]any{
		"patient_id": patientID,
		"items": []map[string]any{
			{"item_type": "service", "description": "Consultation", "quantity": 1, "unit_price": 1000},
		},
	}, &invoice)

	var paid struct {
		Amount float64 `json:"amount"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/billing/invoice/"+invoice.ID.String()+"/payments", token, map[string]any{
		"amount": 400, "method": "cash",
	}, &paid)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, float64(400), paid.Amount)

	var afterPartial struct {
		DueAmount     float64 `json:"due_amount"`
		PaymentStatus string  `json:"payment_status"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/billing/invoice/"+invoice.ID.String(), token, nil, &afterPartial)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, float64(600), afterPartial.DueAmount)
	assert.Equal(t, "partial", afterPartial.PaymentStatus)

	// Overpay: due is 600, try to pay 9999.
	resp = ts.Do(t, http.MethodPost, "/api/billing/invoice/"+invoice.ID.String()+"/payments", token, map[string]any{
		"amount": 9999, "method": "cash",
	}, nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

