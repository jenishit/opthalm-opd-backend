package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

// ─── Inventory Items ─────────────────────────────────────────────

type InventoryRepository struct {
	DB *postgres.DB
}

func NewInventoryRepository(db *postgres.DB) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

var inventoryItemColumns = []string{
	"id", "category", "sku", "name", "brand", "model", "color", "size",
	"cost_price", "selling_price", "quantity_on_hand", "reorder_threshold",
	"unit", "is_active", "created_by", "updated_by", "created_at", "updated_at",
}

func (r *InventoryRepository) Create(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) (*domain.InventoryItem, error) {
	query, args, err := sq.Insert("inventory_items").
