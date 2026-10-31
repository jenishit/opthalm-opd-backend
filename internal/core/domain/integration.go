package domain

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

// Integration tracks a clinic's connection to one third-party provider.
// OAuth access/refresh tokens are stored alongside this in the repository
// layer but deliberately live outside this struct — see
// IntegrationRepository.ListIntegrations, which never selects them.
type Integration struct {
	ID          uuid.UUID
	ClinicID    uuid.UUID
	Provider    string // xero | twilio | google_calendar | essilor_connect
	Connected   bool
	ConnectedAt *time.Time
}

// KnownIntegrationProviders is every provider the API will acknowledge.
// GET /clinic/integrations always reports one entry per provider here, even
// for clinics that have never connected anything.
var KnownIntegrationProviders = []string{"xero", "twilio", "google_calendar", "essilor_connect"}

func IsKnownIntegrationProvider(provider string) bool {
	return slices.Contains(KnownIntegrationProviders, provider)
}
