package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestReports_SalesDailyReflectsSeededInvoice(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Report Patient", "9800000040")

	ts.Do(t, http.MethodPost, "/api/billing/invoice", token, map[string]any{
		"patient_id": patientID,
		"items": []map[string]any{
			{"item_type": "service", "description": "Consultation", "quantity": 1, "unit_price": 1500},
		},
	}, nil)

	today := time.Now().Format("2006-01-02")
	var summary struct {
		InvoiceCount int     `json:"invoice_count"`
		TotalSales   float64 `json:"total_sales"`
		TotalDue     float64 `json:"total_due"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/reports/sales/daily?date="+today, token, nil, &summary)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, summary.InvoiceCount)
	assert.Equal(t, float64(1500), summary.TotalSales)
	assert.Equal(t, float64(1500), summary.TotalDue)
}

func TestReports_PatientDuesOmitsFullyPaidInvoice(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Report Patient 2", "9800000041")

	var invoice struct {
		ID uuid.UUID `json:"id"`
	}
	ts.DoData(t, http.MethodPost, "/api/billing/invoice", token, map[string]any{
		"patient_id": patientID,
		"items": []map[string]any{
			{"item_type": "service", "description": "Consultation", "quantity": 1, "unit_price": 800},
		},
	}, &invoice)

	var dues []struct {
		InvoiceNo string `json:"invoice_no"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/reports/dues", token, nil, &dues)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, dues, 1, "an unpaid invoice should show up in the dues report")

	ts.Do(t, http.MethodPost, "/api/billing/invoice/"+invoice.ID.String()+"/payments", token, map[string]any{
		"amount": 800, "method": "cash",
	}, nil)

	resp = ts.DoData(t, http.MethodGet, "/api/reports/dues", token, nil, &dues)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, dues, 0, "a fully paid invoice should not show up in the dues report")
}

func TestReports_InventoryValuation(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	createInventoryItem(t, ts, token, "RPT-1", 10) // cost 1000, sell 2000, qty 10

	var valuation struct {
		TotalItems     int     `json:"total_items"`
		TotalUnits     int     `json:"total_units"`
		TotalCostValue float64 `json:"total_cost_value"`
		TotalSellValue float64 `json:"total_sell_value"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/reports/inventory/valuation", token, nil, &valuation)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, valuation.TotalItems)
	assert.Equal(t, 10, valuation.TotalUnits)
	assert.Equal(t, float64(10000), valuation.TotalCostValue)
	assert.Equal(t, float64(20000), valuation.TotalSellValue)
}

func TestReports_CSVAndPDFFormats(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Report Patient 3", "9800000042")
	ts.Do(t, http.MethodPost, "/api/billing/invoice", token, map[string]any{
		"patient_id": patientID,
		"items": []map[string]any{
			{"item_type": "service", "description": "Consultation", "quantity": 1, "unit_price": 100},
		},
	}, nil)

	today := time.Now().Format("2006-01-02")

	resp := ts.Do(t, http.MethodGet, "/api/reports/sales/daily?date="+today+"&format=csv", token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))

	resp = ts.Do(t, http.MethodGet, "/api/reports/sales/daily?date="+today+"&format=pdf", token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/pdf", resp.Header.Get("Content-Type"))
}

func TestReports_NonAdminRejected(t *testing.T) {
	ts := testutil.NewTestServer(t)
	doctorToken, _ := ts.Login(t, "reportsuser@test.local", "ROLE_DOCTOR")

	resp := ts.Do(t, http.MethodGet, "/api/reports/dues", doctorToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
