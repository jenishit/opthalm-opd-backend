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
}

func InventoryItemRes(item *domain.InventoryItem) *InventoryItemResponse {
	return &InventoryItemResponse{
		ID:               item.ID,
		Category:         string(item.Category),
		SKU:              item.SKU,
		Name:             item.Name,
		Brand:            item.Brand,
		Model:            item.Model,
		Color:            item.Color,
		Size:             item.Size,
		CostPrice:        item.CostPrice,
		SellingPrice:     item.SellingPrice,
		QuantityOnHand:   item.QuantityOnHand,
		ReorderThreshold: item.ReorderThreshold,
		Unit:             item.Unit,
		IsActive:         item.IsActive,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}

func InventoryItemResList(items []*domain.InventoryItem) []*InventoryItemResponse {
	res := make([]*InventoryItemResponse, 0, len(items))
	for _, item := range items {
		res = append(res, InventoryItemRes(item))
	}
	return res
}

type StockMovementResponse struct {
	ID              uuid.UUID  `json:"id"`
	InventoryItemID uuid.UUID  `json:"inventory_item_id"`
	MovementType    string     `json:"movement_type"`
	Quantity        int        `json:"quantity"`
	ReferenceType   *string    `json:"reference_type"`
	ReferenceID     *uuid.UUID `json:"reference_id"`
	Notes           *string    `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
}

func StockMovementRes(m *domain.StockMovement) *StockMovementResponse {
	return &StockMovementResponse{
		ID:              m.ID,
		InventoryItemID: m.InventoryItemID,
		MovementType:    string(m.MovementType),
		Quantity:        m.Quantity,
		ReferenceType:   m.ReferenceType,
		ReferenceID:     m.ReferenceID,
		Notes:           m.Notes,
		CreatedAt:       m.CreatedAt,
	}
}

func StockMovementResList(ms []*domain.StockMovement) []*StockMovementResponse {
	res := make([]*StockMovementResponse, 0, len(ms))
	for _, m := range ms {
		res = append(res, StockMovementRes(m))
	}
	return res
}

// ─── Vendors ─────────────────────────────────────────────────────

type CreateVendorReq struct {
	Name          string  `json:"name" binding:"required"`
	ContactPerson *string `json:"contact_person"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
}

type UpdateVendorReq struct {
	Name          *string `json:"name"`
	ContactPerson *string `json:"contact_person"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
}

type VendorResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ContactPerson *string   `json:"contact_person"`
	Phone         *string   `json:"phone"`
	Email         *string   `json:"email"`
	Address       *string   `json:"address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
