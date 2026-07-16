package domain

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID           uuid.UUID
	FullName     string
	DOB          string
	Gender       string
	Phone        string
	Occupation   *string
	RegisteredOn time.Time
	CreatedBy    uuid.UUID
	UpdatedBy    uuid.UUID
	Email        *string
	Address      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
