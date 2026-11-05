package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type CreatePatientReq struct {
	FullName           string  `json:"full_name" binding:"required"`
	DOB                string  `json:"dob" binding:"required"`
	Gender             string  `json:"gender"`
	Phone              string  `json:"phone" binding:"required"`
	Occupation         *string `json:"occupation"`
	Email              *string `json:"email"`
	Address            *string `json:"address"`
	ReferredByName     *string `json:"referred_by_name"`
	ReferredByFacility *string `json:"referred_by_facility"`
}

type PatientResponse struct {
	ID                 uuid.UUID `json:"patient_id"`
	FullName           string    `json:"full_name"`
	DOB                string    `json:"dob"`
	Gender             string    `json:"gender"`
	Phone              string    `json:"phone"`
	Occupation         *string   `json:"occupation"`
	ReferredByName     *string   `json:"referred_by_name,omitempty"`
	ReferredByFacility *string   `json:"referred_by_facility,omitempty"`
	RegisteredOn       string    `json:"registered_on"`
	CreatedBy          uuid.UUID `json:""`
	Email              *string   `json:"email"`
	Address            *string   `json:"address"`
	CreatedAt          time.Time `json:"created_at"`
}

func PatientRes(p *domain.Patient) *PatientResponse {
	return &PatientResponse{
		ID:                 p.ID,
		FullName:           p.FullName,
		DOB:                p.DOB,
		Gender:             p.Gender,
		Phone:              p.Phone,
		Occupation:         p.Occupation,
		ReferredByName:     p.ReferredByName,
		ReferredByFacility: p.ReferredByFacility,
		RegisteredOn:       p.RegisteredOn.String(),
		CreatedBy:          p.CreatedBy,
		Email:              p.Email,
		Address:            p.Address,
		CreatedAt:          p.CreatedAt,
	}
}

func PatientResponses(pts []*domain.Patient) []*PatientResponse {
	patients := make([]*PatientResponse, 0, len(pts))

	for _, p := range pts {
		patients = append(patients, PatientRes(p))
	}
	return patients
}
