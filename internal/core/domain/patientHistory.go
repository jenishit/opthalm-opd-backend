package domain

import (
	"time"

	"github.com/google/uuid"
)

// PatientHistory is one append-only entry in a patient's clinical history
// timeline (background/PMH) — never edited in place, only ever added to.
type PatientHistory struct {
	ID                  uuid.UUID
	PatientID           uuid.UUID
	UsesSpectacles      bool
	SystemicMedications *string
	TopicalMedications  *string
	Notes               *string
	ConditionIDs        []uuid.UUID // references the history_conditions catalog
	CreatedBy           uuid.UUID
	CreatedAt           time.Time
}
