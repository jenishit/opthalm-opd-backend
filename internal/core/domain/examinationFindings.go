package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExaminationFindings struct {
	ID                uuid.UUID
	VisitID           uuid.UUID
	IOPOD             *float64
	IOPOS             *float64
	IOPTime           *string
	BP                *string
	Convergence       *string
	CoverTest         *string
	EOM               *string
	AnteriorSegmentOD *string
	AnteriorSegmentOS *string
	FundusOD          *string
	FundusOS          *string
	CreatedBy         uuid.UUID
	UpdatedBy         uuid.UUID
	UpdatedAt         time.Time
}
