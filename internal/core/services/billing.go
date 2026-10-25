package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type InvoiceService struct {
	repo port.InvoiceRepository
}

func NewInvoiceService(r port.InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: r}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, clinicID uuid.UUID, invoice *domain.Invoice, items []*domain.InvoiceItem) (*domain.InvoiceDetails, error) {
	var subtotal, discount float64
	for _, item := range items {
		item.LineTotal = item.UnitPrice*float64(item.Quantity) - item.DiscountAmount
		subtotal += item.UnitPrice * float64(item.Quantity)
		discount += item.DiscountAmount
	}

	invoice.Subtotal = subtotal
	invoice.DiscountAmount += discount
	invoice.TotalAmount = subtotal - invoice.DiscountAmount + invoice.TaxAmount
	invoice.PaidAmount = 0
	invoice.DueAmount = invoice.TotalAmount
	invoice.PaymentStatus = domain.PaymentUnpaid
	if invoice.Status == "" {
		invoice.Status = domain.InvoiceDraft
	}

	return s.repo.CreateInvoice(ctx, clinicID, invoice, items)
}

func (s *InvoiceService) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.InvoiceDetails, error) {
	return s.repo.GetByID(ctx, clinicID, id)
}

func (s *InvoiceService) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.InvoiceDetails, error) {
	return s.repo.List(ctx, clinicID, limit, offset)
}

func (s *InvoiceService) Search(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.InvoiceDetails, error) {
	return s.repo.Search(ctx, clinicID, query, limit)
}

func (s *InvoiceService) UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.InvoiceStatus, updatedBy uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, clinicID, id, status, updatedBy)
}

func (s *InvoiceService) RecordPayment(ctx context.Context, clinicID uuid.UUID, payment *domain.Payment) (*domain.Payment, error) {
	if payment.Amount <= 0 {
		return nil, domain.ErrInsufficientPayment
	}
	return s.repo.RecordPayment(ctx, clinicID, payment)
}
