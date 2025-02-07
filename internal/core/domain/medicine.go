package domain

import (
	"time"

	"github.com/google/uuid"
)

type Medicine struct {
	ID           uuid.UUID
	MedicineName string
	BrandName    *string
	Strength     *string
	Form         *string
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
