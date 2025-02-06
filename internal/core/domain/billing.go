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
