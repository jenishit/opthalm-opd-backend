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
