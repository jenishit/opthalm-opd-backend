package domain

import (
	"time"

	"github.com/google/uuid"
)

type VisitStatus string

const (
	Scheduled VisitStatus = "scheduled"
	Waiting   VisitStatus = "waiting"
	Examining VisitStatus = "examining"
	Completed VisitStatus = "completed"
	Cancelled VisitStatus = "cancelled"
)

type Visit struct {
	ID             uuid.UUID
	PatientID      uuid.UUID
	ExamineBy      uuid.UUID
	Status         VisitStatus
	VisitDate      string
	CheifComplaint string
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type VisitDetails struct {
	ID             uuid.UUID
	PatientID      uuid.UUID
	PatientName    string
	ExamineBy      uuid.UUID
	ExamineByName  string
	Status         VisitStatus
	VisitDate      string
	CheifComplaint string
	CreatedBy      uuid.UUID
	CreatedByName  string
	UpdatedBy      uuid.UUID
	UpdatedByName  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
