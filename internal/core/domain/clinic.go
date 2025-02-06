package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClinicSettings struct {
	ID             uuid.UUID
	ClinicName     string
	Tagline        *string
	Address        *string
	Phone          *string
	Email          *string
	RegistrationNo *string
	ReportFooter   *string
	UpdatedAt      time.Time
	UpdatedBy      uuid.UUID
}
