package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefractionReading struct {
	ID        uuid.UUID
	VisitID   uuid.UUID
	Eye       string // OD | OS
	Stage     string // unaided | old_rx | pinhole | retinoscopy | subjective_acceptance | final_rx
	Sphere    *float64
	Cylinder  *float64
	Axis      *int
	AddPower  *float64
	VA        *string
	PD        *float64
	LensType  *string // single_vision | bifocal | progressive | contact_lens | none
	CreatedBy uuid.UUID
	CreatedAt time.Time
}
