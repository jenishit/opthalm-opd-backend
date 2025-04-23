package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

// ─── Inventory Items ─────────────────────────────────────────────

type InventoryService struct {
	repo port.InventoryItemRepository
}

func NewInventoryService(r port.InventoryItemRepository) *InventoryService {
	return &InventoryService{repo: r}
}

func (s *InventoryService) Create(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) (*domain.InventoryItem, error) {
	return s.repo.Create(ctx, clinicID, item)
}

func (s *InventoryService) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InventoryItem, error) {
	return s.repo.GetByID(ctx, clinicID, id)
}

func (s *InventoryService) GetBySKU(ctx context.Context, clinicID uuid.UUID, sku string) (*domain.InventoryItem, error) {
	return s.repo.GetBySKU(ctx, clinicID, sku)
}

func (s *InventoryService) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InventoryItem, error) {
	return s.repo.List(ctx, clinicID, limit, offset)
}

func (s *InventoryService) Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InventoryItem, error) {
	return s.repo.Search(ctx, clinicID, query, limit)
}

func (s *InventoryService) Update(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) error {
	return s.repo.Update(ctx, clinicID, item)
}

func (s *InventoryService) Delete(ctx context.Context, clinicID, id uuid.UUID) error {
	return s.repo.Delete(ctx, clinicID, id)
}

func (s *InventoryService) LowStock(ctx context.Context, clinicID uuid.UUID) ([]*domain.InventoryItem, error) {
	return s.repo.LowStock(ctx, clinicID)
}

func (s *InventoryService) AddStock(ctx context.Context, clinicID, itemID uuid.UUID, qty int, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error) {
	referenceType := "manual"
	return s.repo.AddStock(ctx, clinicID, itemID, qty, domain.MovementAdjustmentIn, &referenceType, nil, notes, createdBy)
}

func (s *InventoryService) ListMovements(ctx context.Context, clinicID, itemID uuid.UUID, limit, offset int) ([]*domain.StockMovement, error) {
	return s.repo.ListMovements(ctx, clinicID, itemID, limit, offset)
}

// ─── Vendors ─────────────────────────────────────────────────────

type VendorService struct {
	repo port.VendorRepository
}

func NewVendorService(r port.VendorRepository) *VendorService {
	return &VendorService{repo: r}
}

func (s *VendorService) Create(ctx context.Context, clinicID uuid.UUID, v *domain.Vendor) (*domain.Vendor, error) {
	return s.repo.Create(ctx, clinicID, v)
}

func (s *VendorService) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.Vendor, error) {
	return s.repo.GetByID(ctx, clinicID, id)
}

func (s *VendorService) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.Vendor, error) {
	return s.repo.List(ctx, clinicID, limit, offset)
}

func (s *VendorService) Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.Vendor, error) {
	return s.repo.Search(ctx, clinicID, query, limit)
}

func (s *VendorService) Update(ctx context.Context, clinicID uuid.UUID, v *domain.Vendor) error {
	return s.repo.Update(ctx, clinicID, v)
}

func (s *VendorService) Delete(ctx context.Context, clinicID, id uuid.UUID) error {
	return s.repo.Delete(ctx, clinicID, id)
}

// ─── Stock Purchases ─────────────────────────────────────────────

type StockPurchaseService struct {
