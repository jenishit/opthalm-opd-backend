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
