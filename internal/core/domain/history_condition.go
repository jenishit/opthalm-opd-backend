package domain

import (
	"time"

	"github.com/google/uuid"
)

type HistoryCondition struct {
	ID        uuid.UUID
	Name      string
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
