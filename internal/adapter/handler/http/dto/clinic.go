package dto

import (
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ClinicRequest struct {
	ClinicName               string  `json:"clinic_name" binding:"required"`
	Tagline                  string  `json:"tagline"`
	Address                  string  `json:"address"`
	Phone                    string  `json:"phone"`
	Email                    string  `json:"email"`
	RegistrationNo           string  `json:"registration_no" binding:"required"`
	ReportFooter             string  `json:"report_footer"`
	WorkspaceSlug            *string `json:"workspace_slug"`
	DefaultCurrency          *string `json:"default_currency" binding:"omitempty,len=3"`
	Timezone                 *string `json:"timezone"`
	RequireConsentBeforeExam *bool   `json:"require_consent_before_exam"`
	LowStockNotifications    *bool   `json:"low_stock_notifications"`
}

type ClinicResponse struct {
	ID                       uuid.UUID `json:"id"`
	ClinicName               string    `json:"clinic_name"`
	Tagline                  string    `json:"tagline"`
	Address                  string    `json:"address"`
	Phone                    string    `json:"phone"`
	Email                    string    `json:"email"`
	RegistrationNo           string    `json:"registration_no"`
	ReportFooter             string    `json:"report_footer"`
	WorkspaceSlug            string    `json:"workspace_slug"`
	DefaultCurrency          string    `json:"default_currency"`
	Timezone                 string    `json:"timezone"`
	RequireConsentBeforeExam bool      `json:"require_consent_before_exam"`
	LowStockNotifications    bool      `json:"low_stock_notifications"`
	UpdatedBy                uuid.UUID `json:"updated_by"`
}

func ClinicsResponse(c *domain.ClinicSettings) *ClinicResponse {
	return &ClinicResponse{
		ID:                       c.ID,
		ClinicName:               c.ClinicName,
		Tagline:                  deref(c.Tagline),
		Address:                  deref(c.Address),
		Phone:                    deref(c.Phone),
		Email:                    deref(c.Email),
		RegistrationNo:           deref(c.RegistrationNo),
		ReportFooter:             deref(c.ReportFooter),
		WorkspaceSlug:            deref(c.WorkspaceSlug),
		DefaultCurrency:          deref(c.DefaultCurrency),
		Timezone:                 deref(c.Timezone),
		RequireConsentBeforeExam: derefBool(c.RequireConsentBeforeExam),
		LowStockNotifications:    derefBool(c.LowStockNotifications),
		UpdatedBy:                c.UpdatedBy,
	}
}

func ClinicsResponses(c []*domain.ClinicSettings) []*ClinicResponse {
	clinics := make([]*ClinicResponse, 0, len(c))

	for _, clinic := range c {
		clinics = append(clinics, ClinicsResponse(clinic))
	}
	return clinics
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
