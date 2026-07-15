package domain

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID        uuid.UUID
	FullName  string
	DOB       *time.Time
	Gender    string
	Phone     *string
	Email     *string
	Address   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
