package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type UpsertPrivacyConfigReq struct {
	RetentionYears               int  `json:"retention_years" binding:"min=1"`
	AutoAnonymizeAfterRetention  bool `json:"auto_anonymize_after_retention"`
	AllowPatientDataExport       bool `json:"allow_patient_data_export"`
	AllowPatientDeletionRequests bool `json:"allow_patient_deletion_requests"`
	Require2FA                   bool `json:"require_2fa"`
	SessionTimeoutMinutes        *int `json:"session_timeout_minutes" binding:"omitempty,min=1"`
	LoginAlerts                  bool `json:"login_alerts"`
}

type PrivacyConfigResponse struct {
	ClinicID                     uuid.UUID `json:"clinic_id"`
	RetentionYears               int       `json:"retention_years"`
	AutoAnonymizeAfterRetention  bool      `json:"auto_anonymize_after_retention"`
	AllowPatientDataExport       bool      `json:"allow_patient_data_export"`
	AllowPatientDeletionRequests bool      `json:"allow_patient_deletion_requests"`
	Require2FA                   bool      `json:"require_2fa"`
	SessionTimeoutMinutes        *int      `json:"session_timeout_minutes"`
	LoginAlerts                  bool      `json:"login_alerts"`
}

func PrivacyConfigRes(c *domain.PrivacyConfig) *PrivacyConfigResponse {
	if c == nil {
		return nil
	}

	return &PrivacyConfigResponse{
		ClinicID:                     c.ClinicID,
		RetentionYears:               c.RetentionYears,
		AutoAnonymizeAfterRetention:  c.AutoAnonymizeAfterRetention,
		AllowPatientDataExport:       c.AllowPatientDataExport,
		AllowPatientDeletionRequests: c.AllowPatientDeletionRequests,
		Require2FA:                   c.Require2FA,
		SessionTimeoutMinutes:        c.SessionTimeoutMinutes,
		LoginAlerts:                  c.LoginAlerts,
	}
}

type AuditLogEntryResponse struct {
	ID          uuid.UUID `json:"id"`
	ClinicID    uuid.UUID `json:"clinic_id"`
	ActorUserID uuid.UUID `json:"actor_user_id"`
	Action      string    `json:"action"`
	TargetType  *string   `json:"target_type,omitempty"`
	TargetID    *string   `json:"target_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func AuditLogEntryRes(e *domain.AuditLogEntry) *AuditLogEntryResponse {
	if e == nil {
		return nil
	}

	return &AuditLogEntryResponse{
		ID:          e.ID,
		ClinicID:    e.ClinicID,
		ActorUserID: e.ActorUserID,
		Action:      e.Action,
		TargetType:  e.TargetType,
		TargetID:    e.TargetID,
		CreatedAt:   e.CreatedAt,
	}
}

func AuditLogEntryResList(entries []*domain.AuditLogEntry) []*AuditLogEntryResponse {
	rsp := make([]*AuditLogEntryResponse, 0, len(entries))
	for _, e := range entries {
		rsp = append(rsp, AuditLogEntryRes(e))
	}
	return rsp
}
