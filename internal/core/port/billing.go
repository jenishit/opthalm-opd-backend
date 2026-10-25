package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type InvoiceRepository interface {
	// CreateInvoice inserts the invoice header and all items in a single
	// transaction, decrementing inventory stock for any item that carries an
	// InventoryItemID.
	CreateInvoice(ctx context.Context, clinicID uuid.UUID, invoice *domain.Invoice, items []*domain.InvoiceItem) (*domain.InvoiceDetails, error)
	GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InvoiceDetails, error)
	List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InvoiceDetails, error)
	Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InvoiceDetails, error)
	UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.InvoiceStatus, updatedBy uuid.UUID) error
	// RecordPayment inserts a payment row and updates the invoice's
	// paid_amount/due_amount/payment_status atomically.
	RecordPayment(ctx context.Context, clinicID uuid.UUID, payment *domain.Payment) (*domain.Payment, error)
}

type InvoiceService interface {
	CreateInvoice(ctx context.Context, clinicID uuid.UUID, invoice *domain.Invoice, items []*domain.InvoiceItem) (*domain.InvoiceDetails, error)
	GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InvoiceDetails, error)
	List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InvoiceDetails, error)
	Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InvoiceDetails, error)
	UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.InvoiceStatus, updatedBy uuid.UUID) error
	RecordPayment(ctx context.Context, clinicID uuid.UUID, payment *domain.Payment) (*domain.Payment, error)
}
