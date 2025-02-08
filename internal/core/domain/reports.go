package domain

import (
	"time"

	"github.com/google/uuid"
)

type SalesSummary struct {
	Period       string
	InvoiceCount int
	Subtotal     float64
	Discount     float64
	Tax          float64
	TotalSales   float64
	TotalPaid    float64
	TotalDue     float64
}

type PatientDue struct {
	PatientID    uuid.UUID
	PatientName  string
	PatientPhone string
	InvoiceNo    string
	InvoiceID    uuid.UUID
	TotalAmount  float64
	PaidAmount   float64
	DueAmount    float64
	CreatedAt    time.Time
