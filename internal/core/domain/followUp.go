package domain

import (
	"time"

	"github.com/google/uuid"
)

type FollowUp struct {
	ID             uuid.UUID
	VisitID        uuid.UUID
	FollowUpDate   string // YYYY-MM-DD
	Reason         *string
	ReminderSentAt *time.Time
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	UpdatedAt      time.Time
}
