package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type CreateLabJobReq struct {
	InvoiceID            *uuid.UUID `json:"invoice_id"`
	InvoiceItemID        *uuid.UUID `json:"invoice_item_id"`
	PatientID            uuid.UUID  `json:"patient_id" binding:"required"`
	VendorID             *uuid.UUID `json:"vendor_id"`
	JobType              string     `json:"job_type" binding:"required"`
	ExpectedDeliveryDate *string    `json:"expected_delivery_date"`
	AdvancePayment       float64    `json:"advance_payment"`
	Notes                *string    `json:"notes"`
}

type UpdateLabJobStatusReq struct {
	Status string  `json:"status" binding:"required,oneof=in_fitting ready_to_deliver delivered cancelled"`
	Notes  *string `json:"notes"`
}

type LabJobStatusHistoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	ChangedAt time.Time `json:"changed_at"`
	ChangedBy uuid.UUID `json:"changed_by"`
	Notes     *string   `json:"notes"`
}

type LabJobResponse struct {
	ID                   uuid.UUID                     `json:"id"`
	InvoiceID            *uuid.UUID                    `json:"invoice_id"`
	InvoiceItemID        *uuid.UUID                    `json:"invoice_item_id"`
	PatientID            uuid.UUID                     `json:"patient_id"`
	PatientName          string                        `json:"patient_name"`
	VendorID             *uuid.UUID                    `json:"vendor_id"`
	VendorName           *string                       `json:"vendor_name"`
	JobType              string                        `json:"job_type"`
	Status               string                        `json:"status"`
	ExpectedDeliveryDate *time.Time                    `json:"expected_delivery_date"`
	DeliveredAt          *time.Time                    `json:"delivered_at"`
	AdvancePayment       float64                       `json:"advance_payment"`
	Notes                *string                       `json:"notes"`
	CreatedAt            time.Time                     `json:"created_at"`
	UpdatedAt            time.Time                     `json:"updated_at"`
	History              []LabJobStatusHistoryResponse `json:"history,omitempty"`
}

func LabJobRes(d *domain.LabJobDetails) *LabJobResponse {
	res := &LabJobResponse{
		ID:                   d.ID,
		InvoiceID:            d.InvoiceID,
		InvoiceItemID:        d.InvoiceItemID,
		PatientID:            d.PatientID,
		PatientName:          d.PatientName,
