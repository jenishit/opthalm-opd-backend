package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type InventoryItemRepository interface {
	Create(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) (*domain.InventoryItem, error)
	GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InventoryItem, error)
	GetBySKU(ctx context.Context, clinicID uuid.UUID, sku string) (*domain.InventoryItem, error)
	List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InventoryItem, error)
	Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InventoryItem, error)
	Update(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) error
	Delete(ctx context.Context, clinicID, id uuid.UUID) error
	LowStock(ctx context.Context, clinicID uuid.UUID) ([]*domain.InventoryItem, error)
	AddStock(ctx context.Context, clinicID, itemID uuid.UUID, qty int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error)
	DeductStock(ctx context.Context, clinicID, itemID uuid.UUID, qty int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error)
	ListMovements(ctx context.Context, clinicID, itemID uuid.UUID, limit, offset int) ([]*domain.StockMovement, error)
}

type InventoryItemService interface {
	Create(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) (*domain.InventoryItem, error)
	GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InventoryItem, error)
	GetBySKU(ctx context.Context, clinicID uuid.UUID, sku string) (*domain.InventoryItem, error)
	List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InventoryItem, error)
	Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InventoryItem, error)
	Update(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) error
	Delete(ctx context.Context, clinicID, id uuid.UUID) error
	LowStock(ctx context.Context, clinicID uuid.UUID) ([]*domain.InventoryItem, error)
