package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type CreatePatientReq struct {
	FullName   string  `json:"full_name" binding:"required"`
	DOB        string  `json:"dob" binding:"required"`
	Gender     string  `json:"gender"`
	Phone      string  `json:"phone" binding:"required"`
	Occupation *string `json:"occupation"`
	Email      *string `json:"email"`
	Address    *string `json:"address"`
}

type PatientResponse struct {
	ID           uuid.UUID `json:"patient_id"`
	FullName     string    `json:"full_name"`
	DOB          string    `json:"dob"`
	Gender       string    `json:"gender"`
	Phone        string    `json:"phone"`
	Occupation   *string   `json:"occupation"`
	RegisteredOn string    `json:"registered_on"`
	CreatedBy    uuid.UUID `json:""`
	Email        *string   `json:"email"`
	Address      *string   `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
}

func PatientRes(p *domain.Patient) *PatientResponse {
	return &PatientResponse{
		ID:           p.ID,
		FullName:     p.FullName,
		DOB:          p.DOB,
		Gender:       p.Gender,
		Phone:        p.Phone,
		Occupation:   p.Occupation,
		RegisteredOn: p.RegisteredOn.String(),
		CreatedBy:    p.CreatedBy,
		Email:        p.Email,
		Address:      p.Address,
		CreatedAt:    p.CreatedAt,
	}
}

func PatientResponses(pts []*domain.Patient) []*PatientResponse {
	patients := make([]*PatientResponse, 0, len(pts))

	for _, p := range pts {
		patients = append(patients, &PatientResponse{
			ID:           p.ID,
			FullName:     p.FullName,
			DOB:          p.DOB,
			Gender:       p.Gender,
			Phone:        p.Phone,
			Occupation:   p.Occupation,
			RegisteredOn: p.RegisteredOn.String(),
			CreatedBy:    p.CreatedBy,
			Email:        p.Email,
			Address:      p.Address,
			CreatedAt:    p.CreatedAt,
		})
	}
	return patients
}
