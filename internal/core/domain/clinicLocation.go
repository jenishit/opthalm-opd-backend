package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClinicLocation struct {
	ID        uuid.UUID
	ClinicID  uuid.UUID
	Name      string
	Address   *string
	ExamRooms int
	Status    string // active | opening_soon
	CreatedBy uuid.UUID
	CreatedAt time.Time
}
