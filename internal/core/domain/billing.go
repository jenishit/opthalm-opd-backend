package domain

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceDraft     InvoiceStatus = "draft"
	InvoiceFinalized InvoiceStatus = "finalized"
	InvoiceCancelled InvoiceStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentUnpaid  PaymentStatus = "unpaid"
	PaymentPartial PaymentStatus = "partial"
	PaymentPaid    PaymentStatus = "paid"
	PaymentVoid    PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentCash         PaymentMethod = "cash"
	PaymentCard         PaymentMethod = "card"
	PaymentOnline       PaymentMethod = "online"
	PaymentBankTransfer PaymentMethod = "bank_transfer"
)

type ItemType string

const (
	ItemFrame       ItemType = "frame"
	ItemLens        ItemType = "lens"
	ItemCoating     ItemType = "coating"
	ItemContactLens ItemType = "contact_lens"
	ItemService     ItemType = "service"
	ItemOther       ItemType = "other"
)

type Invoice struct {
	ID             uuid.UUID
	InvoiceNo      string
	PatientID      uuid.UUID
	VisitID        *uuid.UUID
	Status         InvoiceStatus
	Subtotal       float64
	DiscountAmount float64
	TaxAmount      float64
	TotalAmount    float64
	PaidAmount     float64
	DueAmount      float64
	PaymentStatus  PaymentStatus
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type InvoiceItem struct {
