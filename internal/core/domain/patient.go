package domain

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID         uuid.UUID
	ClinicID   uuid.UUID
	FullName   string
	DOB        string
	Gender     string
	Phone      string
	Occupation *string
	// ReferredByName/ReferredByFacility record who sent this patient to the
	// clinic: a referring doctor/optometrist's name and/or the hospital or
	// clinic name, captured once at registration. Both are optional — most
	// patients are walk-ins with no referral.
	ReferredByName     *string
	ReferredByFacility *string
	RegisteredOn       time.Time
	CreatedBy          uuid.UUID
	UpdatedBy          uuid.UUID
	Email              *string
	Address            *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
