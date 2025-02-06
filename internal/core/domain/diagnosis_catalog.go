package domain

import (
	"time"

	"github.com/google/uuid"
)

type DiagnosisCatalog struct {
	ID        uuid.UUID
	Icd10Code *string
	Name      string
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
