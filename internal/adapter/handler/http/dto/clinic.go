package dto

import (
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ClinicRequest struct {
	ClinicName     string `json:"clinic_name" binding:"required"`
	Tagline        string `json:"tagline"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	RegistrationNo string `json:"registration_no" binding:"required"`
	ReportFooter   string `json:"report_footer"`
}

type ClinicResponse struct {
	ID             uuid.UUID `json:"id"`
	ClinicName     string    `json:"clinic_name"`
	Tagline        string    `json:"tagline"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	RegistrationNo string    `json:"registration_no"`
	ReportFooter   string    `json:"report_footer"`
	UpdatedBy      uuid.UUID `json:"updated_by"`
}

func ClinicsResponse(c *domain.ClinicSettings) *ClinicResponse {
	return &ClinicResponse{
		ID:             c.ID,
		ClinicName:     c.ClinicName,
		Tagline:        deref(c.Tagline),
		Address:        deref(c.Address),
		Phone:          deref(c.Phone),
		Email:          deref(c.Email),
		RegistrationNo: deref(c.RegistrationNo),
		ReportFooter:   deref(c.ReportFooter),
		UpdatedBy:      c.UpdatedBy,
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
