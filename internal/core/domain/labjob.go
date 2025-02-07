package domain

import (
	"time"

	"github.com/google/uuid"
)

type LabJobStatus string

const (
	LabJobInFitting      LabJobStatus = "in_fitting"
	LabJobReadyToDeliver LabJobStatus = "ready_to_deliver"
	LabJobDelivered      LabJobStatus = "delivered"
	LabJobCancelled      LabJobStatus = "cancelled"
)

type LabJob struct {
	ID                   uuid.UUID
	InvoiceID            *uuid.UUID
	InvoiceItemID        *uuid.UUID
	PatientID            uuid.UUID
	VendorID             *uuid.UUID
	JobType              string
	Status               LabJobStatus
	ExpectedDeliveryDate *time.Time
	DeliveredAt          *time.Time
	AdvancePayment       float64
	Notes                *string
	CreatedBy            uuid.UUID
	UpdatedBy            uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type LabJobStatusHistory struct {
	ID        uuid.UUID
	LabJobID  uuid.UUID
	Status    LabJobStatus
	ChangedAt time.Time
	ChangedBy uuid.UUID
	Notes     *string
}

// LabJobDetails is the denormalized view used for read endpoints.
type LabJobDetails struct {
	LabJob
	PatientName string
	VendorName  *string
	History     []*LabJobStatusHistory
}
