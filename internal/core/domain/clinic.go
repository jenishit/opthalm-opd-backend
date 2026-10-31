package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClinicSettings struct {
	ID                       uuid.UUID
	ClinicName               string
	Tagline                  *string
	Address                  *string
	Phone                    *string
	Email                    *string
	RegistrationNo           *string
	ReportFooter             *string
	WorkspaceSlug            *string
	DefaultCurrency          *string
	Timezone                 *string
	RequireConsentBeforeExam *bool
	LowStockNotifications    *bool
	UpdatedAt                time.Time
	UpdatedBy                uuid.UUID
}
