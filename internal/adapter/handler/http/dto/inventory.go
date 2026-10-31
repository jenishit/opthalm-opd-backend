package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
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
}

func VendorRes(v *domain.Vendor) *VendorResponse {
	return &VendorResponse{
		ID:            v.ID,
		Name:          v.Name,
		ContactPerson: v.ContactPerson,
		Phone:         v.Phone,
		Email:         v.Email,
		Address:       v.Address,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}
}

func VendorResList(vs []*domain.Vendor) []*VendorResponse {
	res := make([]*VendorResponse, 0, len(vs))
	for _, v := range vs {
		res = append(res, VendorRes(v))
	}
	return res
}

// ─── Stock Purchases ─────────────────────────────────────────────

type CreateStockPurchaseItemReq struct {
	InventoryItemID uuid.UUID `json:"inventory_item_id" binding:"required"`
	Quantity        int       `json:"quantity" binding:"required,min=1"`
	UnitCost        float64   `json:"unit_cost" binding:"required,min=0"`
}

type CreateStockPurchaseReq struct {
	VendorID     uuid.UUID                    `json:"vendor_id" binding:"required"`
	InvoiceRefNo *string                      `json:"invoice_ref_no"`
	PaidAmount   float64                      `json:"paid_amount"`
	Items        []CreateStockPurchaseItemReq `json:"items" binding:"required,min=1,dive"`
}

type StockPurchaseItemResponse struct {
	ID              uuid.UUID `json:"id"`
	InventoryItemID uuid.UUID `json:"inventory_item_id"`
	Quantity        int       `json:"quantity"`
	UnitCost        float64   `json:"unit_cost"`
	LineTotal       float64   `json:"line_total"`
}

type StockPurchaseResponse struct {
	ID           uuid.UUID                   `json:"id"`
	VendorID     uuid.UUID                   `json:"vendor_id"`
	VendorName   string                      `json:"vendor_name"`
	PurchaseDate time.Time                   `json:"purchase_date"`
	InvoiceRefNo *string                     `json:"invoice_ref_no"`
	TotalAmount  float64                     `json:"total_amount"`
	PaidAmount   float64                     `json:"paid_amount"`
	DueAmount    float64                     `json:"due_amount"`
	CreatedAt    time.Time                   `json:"created_at"`
	Items        []StockPurchaseItemResponse `json:"items,omitempty"`
}

func StockPurchaseRes(d *domain.StockPurchaseDetails) *StockPurchaseResponse {
	res := &StockPurchaseResponse{
		ID:           d.ID,
		VendorID:     d.VendorID,
		VendorName:   d.VendorName,
		PurchaseDate: d.PurchaseDate,
		InvoiceRefNo: d.InvoiceRefNo,
		TotalAmount:  d.TotalAmount,
		PaidAmount:   d.PaidAmount,
		DueAmount:    d.DueAmount,
		CreatedAt:    d.CreatedAt,
	}
	for _, it := range d.Items {
		res.Items = append(res.Items, StockPurchaseItemResponse{
			ID:              it.ID,
			InventoryItemID: it.InventoryItemID,
			Quantity:        it.Quantity,
			UnitCost:        it.UnitCost,
			LineTotal:       it.LineTotal,
		})
	}
	return res
}

func StockPurchaseResList(ds []*domain.StockPurchaseDetails) []*StockPurchaseResponse {
	res := make([]*StockPurchaseResponse, 0, len(ds))
	for _, d := range ds {
		res = append(res, StockPurchaseRes(d))
	}
	return res
}
