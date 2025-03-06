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

type InvoiceRepository struct {
	DB    *postgres.DB
	Stock *InventoryRepository
}

func NewInvoiceRepository(db *postgres.DB, stock *InventoryRepository) *InvoiceRepository {
	return &InvoiceRepository{DB: db, Stock: stock}
}

func (r *InvoiceRepository) CreateInvoice(ctx context.Context, clinicID uuid.UUID, invoice *domain.Invoice, items []*domain.InvoiceItem) (*domain.InvoiceDetails, error) {
	referenceType := "invoice"

	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		query, args, err := sq.Insert("invoices").
			Columns(
				"clinic_id", "patient_id", "visit_id", "status", "subtotal", "discount_amount",
				"tax_amount", "total_amount", "paid_amount", "due_amount",
				"payment_status", "created_by", "updated_by",
			).
			Values(
				clinicID, invoice.PatientID, nullUUIDPtr(invoice.VisitID), string(invoice.Status), invoice.Subtotal, invoice.DiscountAmount,
				invoice.TaxAmount, invoice.TotalAmount, invoice.PaidAmount, invoice.DueAmount,
				string(invoice.PaymentStatus), invoice.CreatedBy, invoice.UpdatedBy,
			).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("InvoiceRepo.CreateInvoice build invoice: %w", err)
		}

		if err := tx.QueryRow(ctx, query, args...).Scan(&invoice.ID); err != nil {
			return fmt.Errorf("InvoiceRepo.CreateInvoice insert invoice: %w", err)
		}

		for _, item := range items {
			ib := sq.Insert("invoice_items").
				Columns(
					"clinic_id", "invoice_id", "bundle_id", "item_type", "description", "inventory_item_id",
					"quantity", "unit_price", "discount_amount", "line_total",
				).
				Values(
					clinicID, invoice.ID, nullUUIDPtr(item.BundleID), string(item.ItemType), item.Description, nullUUIDPtr(item.InventoryItemID),
					item.Quantity, item.UnitPrice, item.DiscountAmount, item.LineTotal,
				).
				Suffix("RETURNING id").
				PlaceholderFormat(sq.Dollar)

			query, args, err := ib.ToSql()
			if err != nil {
				return fmt.Errorf("InvoiceRepo.CreateInvoice build item: %w", err)
			}
			if err := tx.QueryRow(ctx, query, args...).Scan(&item.ID); err != nil {
