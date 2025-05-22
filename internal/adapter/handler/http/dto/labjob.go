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
