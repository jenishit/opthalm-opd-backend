package domain

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type GetProfileDetails struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     *string   `json:"phone"`
	RoleName  string    `json:"role_name"`
	Email     string    `json:"email"`
}
