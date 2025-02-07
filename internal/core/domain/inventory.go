package domain

import (
	"time"

	"github.com/google/uuid"
)

type InventoryCategory string

const (
	CategoryFrame       InventoryCategory = "frame"
	CategoryLens        InventoryCategory = "lens"
	CategoryContactLens InventoryCategory = "contact_lens"
	CategorySunglasses  InventoryCategory = "sunglasses"
	CategoryCoating     InventoryCategory = "coating"
	CategoryAccessory   InventoryCategory = "accessory"
	CategoryOther       InventoryCategory = "other"
)

type StockMovementType string

const (
	MovementPurchaseIn   StockMovementType = "purchase_in"
	MovementSaleOut      StockMovementType = "sale_out"
	MovementAdjustmentIn StockMovementType = "adjustment_in"
	MovementAdjustOut    StockMovementType = "adjustment_out"
	MovementReturnIn     StockMovementType = "return_in"
)

type InventoryItem struct {
	ID               uuid.UUID
	ClinicID         uuid.UUID
	Category         InventoryCategory
	SKU              string
	Name             string
	Brand            *string
	Model            *string
	Color            *string
	Size             *string
	CostPrice        float64
	SellingPrice     float64
	QuantityOnHand   int
	ReorderThreshold int
	Unit             string
	IsActive         bool
	DeletedAt        *time.Time
	CreatedBy        uuid.UUID
	UpdatedBy        uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type StockMovement struct {
