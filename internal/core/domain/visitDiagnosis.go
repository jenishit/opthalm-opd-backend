package domain

import (
	"time"

	"github.com/google/uuid"
)

type VisitDiagnosis struct {
	ID            uuid.UUID
	VisitID       uuid.UUID
	DiagnosisID   uuid.UUID // catalog id, or a local finding id from the exam chart
	DiagnosisName *string
	Icd10Code     *string
	Laterality    string  // OD | OS | OU | none
	Status        *string // provisional | confirmed | resolved
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
}
