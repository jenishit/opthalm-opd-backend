package domain

import (
	"time"

	"github.com/google/uuid"
)

type Treatment struct {
	ID            uuid.UUID
	VisitID       uuid.UUID
	Type          string // medicine | procedure | advice
	MedicineID    *uuid.UUID
	MedicineName  *string
	ProcedureName *string
	AdviceText    *string
	Route         *string
	Dosage        *string
	Frequency     *string
	Duration      *string
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
}
