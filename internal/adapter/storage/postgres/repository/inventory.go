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
		Columns(
			"clinic_id", "category", "sku", "name", "brand", "model", "color", "size",
			"cost_price", "selling_price", "quantity_on_hand", "reorder_threshold",
			"unit", "created_by", "updated_by",
		).
		Values(
			clinicID, string(item.Category), item.SKU, item.Name, nullStringPtr(item.Brand), nullStringPtr(item.Model), nullStringPtr(item.Color), nullStringPtr(item.Size),
			item.CostPrice, item.SellingPrice, item.QuantityOnHand, item.ReorderThreshold,
			item.Unit, item.CreatedBy, item.UpdatedBy,
		).
		Suffix("RETURNING id, is_active, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InventoryRepo.Create build: %w", err)
	}

	if err := r.DB.QueryRow(ctx, query, args...).Scan(&item.ID, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
		return nil, fmt.Errorf("InventoryRepo.Create exec: %w", err)
	}
	item.ClinicID = clinicID
	return item, nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InventoryItem, error) {
	qb := sq.Select(inventoryItemColumns...).
		From("inventory_items").
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)
	return scanInventoryItem(ctx, r.DB, qb)
