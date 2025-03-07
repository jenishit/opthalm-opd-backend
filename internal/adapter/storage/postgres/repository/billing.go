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
				return fmt.Errorf("InvoiceRepo.CreateInvoice insert item: %w", err)
			}
			item.InvoiceID = invoice.ID

			if item.InventoryItemID != nil {
				itemID := item.ID
				if _, err := r.Stock.deductStockTx(ctx, tx, clinicID, *item.InventoryItemID, item.Quantity, domain.MovementSaleOut, &referenceType, &itemID, nil, invoice.CreatedBy); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, clinicID, invoice.ID)
}

func (r *InvoiceRepository) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InvoiceDetails, error) {
	query, args, err := invoiceSelect().
		Where(sq.Eq{"i.id": id, "i.clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InvoiceRepo.GetByID build: %w", err)
	}

	details, err := scanInvoiceDetails(ctx, r.DB, query, args)
	if err != nil {
		return nil, err
	}

	if details.Items, err = r.listItems(ctx, id); err != nil {
		return nil, err
	}
	if details.Payments, err = r.listPayments(ctx, id); err != nil {
		return nil, err
	}

	return details, nil
}

func (r *InvoiceRepository) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InvoiceDetails, error) {
	query, args, err := invoiceSelect().
		Where(sq.Eq{"i.clinic_id": clinicID}).
		OrderBy("i.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InvoiceRepo.List build: %w", err)
	}
	return scanInvoiceDetailsList(ctx, r.DB, query, args)
}

func (r *InvoiceRepository) Search(ctx context.Context, clinicID uuid.UUID, queryStr string, limit int) ([]*domain.InvoiceDetails, error) {
	query, args, err := invoiceSelect().
		Where(sq.Eq{"i.clinic_id": clinicID}).
		Where(sq.Or{
			sq.Expr("i.invoice_no ILIKE '%' || ? || '%'", queryStr),
			sq.Expr("p.full_name ILIKE '%' || ? || '%'", queryStr),
			sq.Expr("p.phone ILIKE '%' || ? || '%'", queryStr),
		}).
		OrderBy("i.created_at DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InvoiceRepo.Search build: %w", err)
	}
	return scanInvoiceDetailsList(ctx, r.DB, query, args)
}

func (r *InvoiceRepository) UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.InvoiceStatus, updatedBy uuid.UUID) error {
	query, args, err := sq.Update("invoices").
		Set("status", string(status)).
		Set("updated_by", updatedBy).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("InvoiceRepo.UpdateStatus build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("InvoiceRepo.UpdateStatus exec: %w", err)
	}
	return nil
}

func (r *InvoiceRepository) RecordPayment(ctx context.Context, clinicID uuid.UUID, payment *domain.Payment) (*domain.Payment, error) {
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		var totalAmount, paidAmount float64
		var paymentStatus string
		err := tx.QueryRow(ctx,
			`SELECT total_amount, paid_amount, payment_status FROM invoices WHERE id = $1 AND clinic_id = $2 FOR UPDATE`,
			payment.InvoiceID, clinicID,
		).Scan(&totalAmount, &paidAmount, &paymentStatus)
		if err != nil {
			if err == pgx.ErrNoRows {
				return domain.ErrDataNotFound
			}
			return fmt.Errorf("InvoiceRepo.RecordPayment lock invoice: %w", err)
		}

		if domain.PaymentStatus(paymentStatus) == domain.PaymentPaid {
			return domain.ErrInvoiceAlreadyPaid
		}

		newPaid := paidAmount + payment.Amount
		newDue := totalAmount - newPaid
		if newDue < 0 {
			return domain.ErrPaymentExceedsDueAmount
		}

		newStatus := domain.PaymentPartial
		if newDue == 0 {
			newStatus = domain.PaymentPaid
		} else if newPaid == 0 {
			newStatus = domain.PaymentUnpaid
		}

		query, args, err := sq.Insert("payments").
			Columns("clinic_id", "invoice_id", "amount", "method", "reference_no", "created_by").
			Values(clinicID, payment.InvoiceID, payment.Amount, string(payment.Method), nullStringPtr(payment.ReferenceNo), payment.CreatedBy).
			Suffix("RETURNING id, paid_at, created_at").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("InvoiceRepo.RecordPayment build payment: %w", err)
		}
		if err := tx.QueryRow(ctx, query, args...).Scan(&payment.ID, &payment.PaidAt, &payment.CreatedAt); err != nil {
			return fmt.Errorf("InvoiceRepo.RecordPayment insert payment: %w", err)
		}

		uQuery, uArgs, err := sq.Update("invoices").
			Set("paid_amount", newPaid).
			Set("due_amount", newDue).
			Set("payment_status", string(newStatus)).
			Set("updated_at", sq.Expr("NOW()")).
			Where(sq.Eq{"id": payment.InvoiceID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("InvoiceRepo.RecordPayment build invoice update: %w", err)
		}
		if _, err := tx.Exec(ctx, uQuery, uArgs...); err != nil {
			return fmt.Errorf("InvoiceRepo.RecordPayment update invoice: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *InvoiceRepository) listItems(ctx context.Context, invoiceID uuid.UUID) ([]*domain.InvoiceItem, error) {
	query, args, err := sq.Select(
		"id", "invoice_id", "bundle_id", "item_type", "description", "inventory_item_id",
		"quantity", "unit_price", "discount_amount", "line_total", "created_at",
	).
