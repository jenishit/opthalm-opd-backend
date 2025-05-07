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
