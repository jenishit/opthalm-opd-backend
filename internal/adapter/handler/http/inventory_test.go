package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestInventory_CreateSearchGetLowStock(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	itemID := createInventoryItem(t, ts, token, "INV-1", 2) // reorder_threshold=3 in helper -> low stock

	var got struct {
		SKU            string `json:"sku"`
		QuantityOnHand int    `json:"quantity_on_hand"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/inventory/items/"+itemID.String(), token, nil, &got)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "INV-1", got.SKU)

	var search []struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/inventory/items/search?query=Test+Frame", token, nil, &search)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, search, 1)

	var byBarcode struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/inventory/items/barcode/INV-1", token, nil, &byBarcode)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, itemID, byBarcode.ID)

	var lowStock []struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/inventory/low-stock", token, nil, &lowStock)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, lowStock, 1, "quantity 2 <= reorder_threshold 3 should show up as low stock")
	assert.Equal(t, itemID, lowStock[0].ID)

	resp = ts.Do(t, http.MethodGet, "/api/inventory/items/"+itemID.String()+"/barcode-image", token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
}

func TestInventory_StockPurchaseIncrementsAndWritesMovement(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	itemID := createInventoryItem(t, ts, token, "INV-2", 10)

	var vendor struct {
		ID uuid.UUID `json:"id"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/inventory/vendors", token, map[string]any{
		"name": "Test Vendor",
	}, &vendor)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodPost, "/api/inventory/purchases", token, map[string]any{
		"vendor_id":   vendor.ID,
		"paid_amount": 5000,
		"items": []map[string]any{
			{"inventory_item_id": itemID, "quantity": 20, "unit_cost": 1000},
		},
