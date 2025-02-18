package domain

import (
	"time"

	"github.com/google/uuid"
)

type VisitSymptom struct {
	VisitID                  uuid.UUID
	ForeignBodySensation     *string
	Redness                  *string
	Swelling                 *string
	DiminishedDistanceVision *string
	DiminishedNearVision     *string
	Blurring                 *string
	Headache                 *bool
	EyePain                  *string
	WateryEyes               *string
	Discharge                *string
	Itching                  *string
	PhotoPhobia              *string
	Floaters                 *string
	Flashes                  *string
	DoubleVision             *bool
	RecentFever              *bool
	ComplaintDuration        *string
	CreatedBy                uuid.UUID
	UpdatedBy                uuid.UUID
	UpdatedAt                time.Time
}

type CreateVSResponse struct {
	VisitID                  uuid.UUID
	CreatedBy                uuid.UUID
	UpdatedBy                uuid.UUID
	UpdatedAt                time.Time
}