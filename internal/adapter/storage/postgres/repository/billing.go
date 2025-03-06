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
