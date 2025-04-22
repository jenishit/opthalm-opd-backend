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

