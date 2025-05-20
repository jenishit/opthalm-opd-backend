package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

// ─── Inventory Items ─────────────────────────────────────────────

type CreateInventoryItemReq struct {
	Category         string  `json:"category" binding:"required,oneof=frame lens contact_lens sunglasses coating accessory other"`
	SKU              string  `json:"sku" binding:"required"`
	Name             string  `json:"name" binding:"required"`
	Brand            *string `json:"brand"`
	Model            *string `json:"model"`
	Color            *string `json:"color"`
	Size             *string `json:"size"`
	CostPrice        float64 `json:"cost_price" binding:"min=0"`
	SellingPrice     float64 `json:"selling_price" binding:"min=0"`
	QuantityOnHand   int     `json:"quantity_on_hand" binding:"min=0"`
	ReorderThreshold int     `json:"reorder_threshold" binding:"min=0"`
	Unit             string  `json:"unit"`
}

type UpdateInventoryItemReq struct {
	Category         *string  `json:"category" binding:"omitempty,oneof=frame lens contact_lens sunglasses coating accessory other"`
	Name             *string  `json:"name"`
	Brand            *string  `json:"brand"`
	Model            *string  `json:"model"`
	Color            *string  `json:"color"`
	Size             *string  `json:"size"`
	CostPrice        *float64 `json:"cost_price"`
	SellingPrice     *float64 `json:"selling_price"`
	ReorderThreshold *int     `json:"reorder_threshold"`
}

type AddStockReq struct {
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Notes    *string `json:"notes"`
}

type InventoryItemResponse struct {
	ID               uuid.UUID `json:"id"`
	Category         string    `json:"category"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	Brand            *string   `json:"brand"`
	Model            *string   `json:"model"`
	Color            *string   `json:"color"`
	Size             *string   `json:"size"`
	CostPrice        float64   `json:"cost_price"`
	SellingPrice     float64   `json:"selling_price"`
	QuantityOnHand   int       `json:"quantity_on_hand"`
	ReorderThreshold int       `json:"reorder_threshold"`
	Unit             string    `json:"unit"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
