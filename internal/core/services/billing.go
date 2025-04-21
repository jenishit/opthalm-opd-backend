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
