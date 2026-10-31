package domain

import (
	"time"

	"github.com/google/uuid"
)

type PrivacyConfig struct {
	ClinicID                     uuid.UUID
	RetentionYears               int
	AutoAnonymizeAfterRetention  bool
	AllowPatientDataExport       bool
	AllowPatientDeletionRequests bool
	Require2FA                   bool
	SessionTimeoutMinutes        *int
	LoginAlerts                  bool
}

type AuditLogEntry struct {
	ID          uuid.UUID
	ClinicID    uuid.UUID
	ActorUserID uuid.UUID
	Action      string
	TargetType  *string
	TargetID    *string
	CreatedAt   time.Time
}
