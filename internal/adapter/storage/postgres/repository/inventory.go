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
}

func (r *InventoryRepository) GetBySKU(ctx context.Context, clinicID uuid.UUID, sku string) (*domain.InventoryItem, error) {
	qb := sq.Select(inventoryItemColumns...).
		From("inventory_items").
		Where(sq.Eq{"sku": sku, "clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)
	return scanInventoryItem(ctx, r.DB, qb)
}

func (r *InventoryRepository) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InventoryItem, error) {
	qb := sq.Select(inventoryItemColumns...).
		From("inventory_items").
		Where(sq.Eq{"clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)
	return scanInventoryItems(ctx, r.DB, qb)
}

func (r *InventoryRepository) Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InventoryItem, error) {
	qb := sq.Select(inventoryItemColumns...).
		From("inventory_items").
		Where(sq.Eq{"clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		Where(sq.Or{
			sq.Expr("name ILIKE '%' || ? || '%'", query),
			sq.Expr("sku ILIKE '%' || ? || '%'", query),
			sq.Expr("brand ILIKE '%' || ? || '%'", query),
		}).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)
	return scanInventoryItems(ctx, r.DB, qb)
}

func (r *InventoryRepository) Update(ctx context.Context, clinicID uuid.UUID, item *domain.InventoryItem) error {
	query, args, err := sq.Update("inventory_items").
		Set("category", sq.Expr("COALESCE(?, category)", nullString(string(item.Category)))).
		Set("name", sq.Expr("COALESCE(?, name)", nullString(item.Name))).
		Set("brand", sq.Expr("COALESCE(?, brand)", nullStringPtr(item.Brand))).
		Set("model", sq.Expr("COALESCE(?, model)", nullStringPtr(item.Model))).
		Set("color", sq.Expr("COALESCE(?, color)", nullStringPtr(item.Color))).
		Set("size", sq.Expr("COALESCE(?, size)", nullStringPtr(item.Size))).
		Set("cost_price", item.CostPrice).
		Set("selling_price", item.SellingPrice).
		Set("reorder_threshold", item.ReorderThreshold).
		Set("updated_by", item.UpdatedBy).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": item.ID, "clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("InventoryRepo.Update build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("InventoryRepo.Update exec: %w", err)
	}
	return nil
}
