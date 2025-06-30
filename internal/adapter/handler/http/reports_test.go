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
