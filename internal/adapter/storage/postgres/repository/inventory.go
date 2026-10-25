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

func (r *InventoryRepository) Delete(ctx context.Context, clinicID, id uuid.UUID) error {
	query, args, err := sq.Update("inventory_items").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("InventoryRepo.Delete build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("InventoryRepo.Delete exec: %w", err)
	}
	return nil
}

func (r *InventoryRepository) LowStock(ctx context.Context, clinicID uuid.UUID) ([]*domain.InventoryItem, error) {
	qb := sq.Select(inventoryItemColumns...).
		From("inventory_items").
		Where(sq.Eq{"clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		Where("is_active = TRUE").
		Where("quantity_on_hand <= reorder_threshold").
		OrderBy("quantity_on_hand ASC").
		PlaceholderFormat(sq.Dollar)
	return scanInventoryItems(ctx, r.DB, qb)
}

func (r *InventoryRepository) AddStock(ctx context.Context, clinicID, itemID uuid.UUID, qty int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error) {
	var movement *domain.StockMovement
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		m, err := r.addStockTx(ctx, tx, clinicID, itemID, qty, movementType, referenceType, referenceID, notes, createdBy)
		movement = m
		return err
	})
	if err != nil {
		return nil, err
	}
	return movement, nil
}

func (r *InventoryRepository) DeductStock(ctx context.Context, clinicID, itemID uuid.UUID, qty int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error) {
	var movement *domain.StockMovement
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		m, err := r.deductStockTx(ctx, tx, clinicID, itemID, qty, movementType, referenceType, referenceID, notes, createdBy)
		movement = m
		return err
	})
	if err != nil {
		return nil, err
	}
	return movement, nil
}

// addStockTx and deductStockTx do the actual locked read + movement insert +
// quantity update within an already-open transaction, so other repositories
// (e.g. StockPurchaseRepository, InvoiceRepository) can compose stock
// changes atomically with their own writes by sharing the same pgx.Tx.
func (r *InventoryRepository) addStockTx(ctx context.Context, tx pgx.Tx, clinicID, itemID uuid.UUID, qty int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error) {
	return r.applyStockDeltaTx(ctx, tx, clinicID, itemID, qty, movementType, referenceType, referenceID, notes, createdBy)
}

func (r *InventoryRepository) deductStockTx(ctx context.Context, tx pgx.Tx, clinicID, itemID uuid.UUID, qty int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error) {
	return r.applyStockDeltaTx(ctx, tx, clinicID, itemID, -qty, movementType, referenceType, referenceID, notes, createdBy)
}

// applyStockDeltaTx locks the inventory_items row, applies delta (positive to
// add, negative to deduct) to quantity_on_hand, and records the movement.
// The movement's own Quantity is always stored positive (direction is
// carried by movementType); delta's sign only decides add vs. deduct.
func (r *InventoryRepository) applyStockDeltaTx(ctx context.Context, tx pgx.Tx, clinicID, itemID uuid.UUID, delta int, movementType domain.StockMovementType, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy uuid.UUID) (*domain.StockMovement, error) {
	var currentQty int
	err := tx.QueryRow(ctx, `SELECT quantity_on_hand FROM inventory_items WHERE id = $1 AND clinic_id = $2 AND deleted_at IS NULL FOR UPDATE`, itemID, clinicID).Scan(&currentQty)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("InventoryRepo lock item: %w", err)
	}

	newQty := currentQty + delta
	if newQty < 0 {
		return nil, domain.ErrInsufficientStock
	}

	if _, err := tx.Exec(ctx, `UPDATE inventory_items SET quantity_on_hand = $1, updated_at = NOW() WHERE id = $2`, newQty, itemID); err != nil {
		return nil, fmt.Errorf("InventoryRepo update quantity: %w", err)
	}

	absQty := delta
	if absQty < 0 {
		absQty = -absQty
	}

	query, args, err := sq.Insert("inventory_stock_movements").
		Columns("clinic_id", "inventory_item_id", "movement_type", "quantity", "reference_type", "reference_id", "notes", "created_by").
		Values(clinicID, itemID, string(movementType), absQty, nullStringPtr(referenceType), nullUUIDPtr(referenceID), nullStringPtr(notes), createdBy).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InventoryRepo build movement: %w", err)
	}

	movement := &domain.StockMovement{
		InventoryItemID: itemID,
		MovementType:    movementType,
		Quantity:        absQty,
		ReferenceType:   referenceType,
		ReferenceID:     referenceID,
		Notes:           notes,
		CreatedBy:       createdBy,
	}
	if err := tx.QueryRow(ctx, query, args...).Scan(&movement.ID, &movement.CreatedAt); err != nil {
		return nil, fmt.Errorf("InventoryRepo insert movement: %w", err)
	}

	return movement, nil
}

func (r *InventoryRepository) ListMovements(ctx context.Context, clinicID, itemID uuid.UUID, limit, offset int) ([]*domain.StockMovement, error) {
	query, args, err := sq.Select(
		"id", "inventory_item_id", "movement_type", "quantity", "reference_type", "reference_id", "notes", "created_by", "created_at",
	).
		From("inventory_stock_movements").
		Where(sq.Eq{"inventory_item_id": itemID, "clinic_id": clinicID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InventoryRepo.ListMovements build: %w", err)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("InventoryRepo.ListMovements query: %w", err)
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		var movementType string
		var referenceType sql.NullString
		var referenceID uuid.NullUUID
		var notes sql.NullString
		if err := rows.Scan(&m.ID, &m.InventoryItemID, &movementType, &m.Quantity, &referenceType, &referenceID, &notes, &m.CreatedBy, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("InventoryRepo.ListMovements scan: %w", err)
		}
		m.MovementType = domain.StockMovementType(movementType)
		if referenceType.Valid {
			m.ReferenceType = &referenceType.String
		}
		if referenceID.Valid {
			m.ReferenceID = &referenceID.UUID
		}
		if notes.Valid {
			m.Notes = &notes.String
		}
		movements = append(movements, &m)
	}
	return movements, rows.Err()
}

func scanInventoryItem(ctx context.Context, db *postgres.DB, qb sq.SelectBuilder) (*domain.InventoryItem, error) {
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	return scanInventoryItemRow(db.QueryRow(ctx, query, args...))
}

func scanInventoryItemRow(row pgx.Row) (*domain.InventoryItem, error) {
	var item domain.InventoryItem
	var category string
	var brand, model, color, size sql.NullString

	err := row.Scan(
		&item.ID, &category, &item.SKU, &item.Name, &brand, &model, &color, &size,
		&item.CostPrice, &item.SellingPrice, &item.QuantityOnHand, &item.ReorderThreshold,
		&item.Unit, &item.IsActive, &item.CreatedBy, &item.UpdatedBy, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("scan inventory item: %w", err)
	}

	item.Category = domain.InventoryCategory(category)
	if brand.Valid {
		item.Brand = &brand.String
	}
	if model.Valid {
		item.Model = &model.String
	}
	if color.Valid {
		item.Color = &color.String
	}
	if size.Valid {
		item.Size = &size.String
	}

	return &item, nil
}

func scanInventoryItems(ctx context.Context, db *postgres.DB, qb sq.SelectBuilder) ([]*domain.InventoryItem, error) {
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query inventory items: %w", err)
	}
	defer rows.Close()

	var items []*domain.InventoryItem
	for rows.Next() {
		item, err := scanInventoryItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ─── Vendors ─────────────────────────────────────────────────────

type VendorRepository struct {
	DB *postgres.DB
}

func NewVendorRepository(db *postgres.DB) *VendorRepository {
	return &VendorRepository{DB: db}
}

var vendorColumns = []string{"id", "name", "contact_person", "phone", "email", "address", "created_at", "updated_at"}

func (r *VendorRepository) Create(ctx context.Context, clinicID uuid.UUID, v *domain.Vendor) (*domain.Vendor, error) {
	query, args, err := sq.Insert("vendors").
		Columns("clinic_id", "name", "contact_person", "phone", "email", "address").
		Values(clinicID, v.Name, nullStringPtr(v.ContactPerson), nullStringPtr(v.Phone), nullStringPtr(v.Email), nullStringPtr(v.Address)).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("VendorRepo.Create build: %w", err)
	}
	if err := r.DB.QueryRow(ctx, query, args...).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, fmt.Errorf("VendorRepo.Create exec: %w", err)
	}
	v.ClinicID = clinicID
	return v, nil
}

func (r *VendorRepository) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.Vendor, error) {
	qb := sq.Select(vendorColumns...).From("vendors").Where(sq.Eq{"id": id, "clinic_id": clinicID}).Where("deleted_at IS NULL").PlaceholderFormat(sq.Dollar)
	return scanVendor(ctx, r.DB, qb)
}

func (r *VendorRepository) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.Vendor, error) {
	qb := sq.Select(vendorColumns...).From("vendors").Where(sq.Eq{"clinic_id": clinicID}).Where("deleted_at IS NULL").
		OrderBy("created_at DESC").Limit(uint64(limit)).Offset(uint64(offset)).PlaceholderFormat(sq.Dollar)
	return scanVendors(ctx, r.DB, qb)
}

func (r *VendorRepository) Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.Vendor, error) {
	qb := sq.Select(vendorColumns...).From("vendors").Where(sq.Eq{"clinic_id": clinicID}).Where("deleted_at IS NULL").
		Where(sq.Or{
			sq.Expr("name ILIKE '%' || ? || '%'", query),
			sq.Expr("phone ILIKE '%' || ? || '%'", query),
		}).Limit(uint64(limit)).PlaceholderFormat(sq.Dollar)
	return scanVendors(ctx, r.DB, qb)
}

func (r *VendorRepository) Update(ctx context.Context, clinicID uuid.UUID, v *domain.Vendor) error {
	query, args, err := sq.Update("vendors").
		Set("name", sq.Expr("COALESCE(?, name)", nullString(v.Name))).
		Set("contact_person", sq.Expr("COALESCE(?, contact_person)", nullStringPtr(v.ContactPerson))).
		Set("phone", sq.Expr("COALESCE(?, phone)", nullStringPtr(v.Phone))).
		Set("email", sq.Expr("COALESCE(?, email)", nullStringPtr(v.Email))).
		Set("address", sq.Expr("COALESCE(?, address)", nullStringPtr(v.Address))).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": v.ID, "clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("VendorRepo.Update build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("VendorRepo.Update exec: %w", err)
	}
	return nil
}

func (r *VendorRepository) Delete(ctx context.Context, clinicID, id uuid.UUID) error {
	query, args, err := sq.Update("vendors").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("VendorRepo.Delete build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("VendorRepo.Delete exec: %w", err)
	}
	return nil
}

func scanVendor(ctx context.Context, db *postgres.DB, qb sq.SelectBuilder) (*domain.Vendor, error) {
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	return scanVendorRow(db.QueryRow(ctx, query, args...))
}

func scanVendorRow(row pgx.Row) (*domain.Vendor, error) {
	var v domain.Vendor
	var contactPerson, phone, email, address sql.NullString
	err := row.Scan(&v.ID, &v.Name, &contactPerson, &phone, &email, &address, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("scan vendor: %w", err)
	}
	if contactPerson.Valid {
		v.ContactPerson = &contactPerson.String
	}
	if phone.Valid {
		v.Phone = &phone.String
	}
	if email.Valid {
		v.Email = &email.String
	}
	if address.Valid {
		v.Address = &address.String
	}
	return &v, nil
}

func scanVendors(ctx context.Context, db *postgres.DB, qb sq.SelectBuilder) ([]*domain.Vendor, error) {
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query vendors: %w", err)
	}
	defer rows.Close()

	var vendors []*domain.Vendor
	for rows.Next() {
		v, err := scanVendorRow(rows)
		if err != nil {
			return nil, err
		}
		vendors = append(vendors, v)
	}
	return vendors, rows.Err()
}

// ─── Stock Purchases ─────────────────────────────────────────────

type StockPurchaseRepository struct {
	DB    *postgres.DB
	Stock *InventoryRepository
}

func NewStockPurchaseRepository(db *postgres.DB, stock *InventoryRepository) *StockPurchaseRepository {
	return &StockPurchaseRepository{DB: db, Stock: stock}
}

func (r *StockPurchaseRepository) CreatePurchase(ctx context.Context, clinicID uuid.UUID, purchase *domain.StockPurchase, items []*domain.StockPurchaseItem) (*domain.StockPurchaseDetails, error) {
	referenceType := "stock_purchase"

	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		query, args, err := sq.Insert("stock_purchases").
			Columns("clinic_id", "vendor_id", "purchase_date", "invoice_ref_no", "total_amount", "paid_amount", "due_amount", "created_by").
			Values(clinicID, purchase.VendorID, purchase.PurchaseDate, nullStringPtr(purchase.InvoiceRefNo), purchase.TotalAmount, purchase.PaidAmount, purchase.DueAmount, purchase.CreatedBy).
			Suffix("RETURNING id, created_at, updated_at").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("StockPurchaseRepo.CreatePurchase build purchase: %w", err)
		}
		if err := tx.QueryRow(ctx, query, args...).Scan(&purchase.ID, &purchase.CreatedAt, &purchase.UpdatedAt); err != nil {
			return fmt.Errorf("StockPurchaseRepo.CreatePurchase insert purchase: %w", err)
		}

		for _, item := range items {
			iQuery, iArgs, err := sq.Insert("stock_purchase_items").
				Columns("clinic_id", "purchase_id", "inventory_item_id", "quantity", "unit_cost", "line_total").
				Values(clinicID, purchase.ID, item.InventoryItemID, item.Quantity, item.UnitCost, item.LineTotal).
				Suffix("RETURNING id").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				return fmt.Errorf("StockPurchaseRepo.CreatePurchase build item: %w", err)
			}
			if err := tx.QueryRow(ctx, iQuery, iArgs...).Scan(&item.ID); err != nil {
				return fmt.Errorf("StockPurchaseRepo.CreatePurchase insert item: %w", err)
			}
			item.PurchaseID = purchase.ID

			purchaseID := purchase.ID
			if _, err := r.Stock.addStockTx(ctx, tx, clinicID, item.InventoryItemID, item.Quantity, domain.MovementPurchaseIn, &referenceType, &purchaseID, nil, purchase.CreatedBy); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, clinicID, purchase.ID)
}

func (r *StockPurchaseRepository) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.StockPurchaseDetails, error) {
	query, args, err := stockPurchaseSelect().
		Where(sq.Eq{"sp.id": id, "sp.clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StockPurchaseRepo.GetByID build: %w", err)
	}

	details, err := scanStockPurchaseDetailsRow(r.DB.QueryRow(ctx, query, args...))
	if err != nil {
		return nil, err
	}

	items, err := r.listItems(ctx, id)
	if err != nil {
		return nil, err
	}
	details.Items = items

	return details, nil
}

func (r *StockPurchaseRepository) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.StockPurchaseDetails, error) {
	query, args, err := stockPurchaseSelect().
		Where(sq.Eq{"sp.clinic_id": clinicID}).
		OrderBy("sp.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StockPurchaseRepo.List build: %w", err)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("StockPurchaseRepo.List query: %w", err)
	}
	defer rows.Close()

	var result []*domain.StockPurchaseDetails
	for rows.Next() {
		d, err := scanStockPurchaseDetailsRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func (r *StockPurchaseRepository) listItems(ctx context.Context, purchaseID uuid.UUID) ([]*domain.StockPurchaseItem, error) {
	query, args, err := sq.Select("id", "purchase_id", "inventory_item_id", "quantity", "unit_cost", "line_total").
		From("stock_purchase_items").
		Where(sq.Eq{"purchase_id": purchaseID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StockPurchaseRepo.listItems build: %w", err)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("StockPurchaseRepo.listItems query: %w", err)
	}
	defer rows.Close()

	var items []*domain.StockPurchaseItem
	for rows.Next() {
		var it domain.StockPurchaseItem
		if err := rows.Scan(&it.ID, &it.PurchaseID, &it.InventoryItemID, &it.Quantity, &it.UnitCost, &it.LineTotal); err != nil {
			return nil, fmt.Errorf("StockPurchaseRepo.listItems scan: %w", err)
		}
		items = append(items, &it)
	}
	return items, rows.Err()
}

func stockPurchaseSelect() sq.SelectBuilder {
	return sq.Select(
		"sp.id", "sp.vendor_id", "v.name", "sp.purchase_date", "sp.invoice_ref_no",
		"sp.total_amount", "sp.paid_amount", "sp.due_amount", "sp.created_by", "sp.created_at", "sp.updated_at",
	).
		From("stock_purchases sp").
		Join("vendors v ON v.id = sp.vendor_id")
}

func scanStockPurchaseDetailsRow(row pgx.Row) (*domain.StockPurchaseDetails, error) {
	var d domain.StockPurchaseDetails
	var invoiceRefNo sql.NullString

	err := row.Scan(
		&d.ID, &d.VendorID, &d.VendorName, &d.PurchaseDate, &invoiceRefNo,
		&d.TotalAmount, &d.PaidAmount, &d.DueAmount, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("scan stock purchase: %w", err)
	}
	if invoiceRefNo.Valid {
		d.InvoiceRefNo = &invoiceRefNo.String
	}
	return &d, nil
}
