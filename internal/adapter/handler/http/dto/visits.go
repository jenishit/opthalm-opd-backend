package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type EntityRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PatientRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type VisitResponse struct {
	ID             uuid.UUID          `json:"id"`
	Patient        PatientRef         `json:"patient"`
	ExamineBy      EntityRef          `json:"examine_by"`
	Status         domain.VisitStatus `json:"status"`
	VisitDate      string             `json:"visit_date"`
	CheifComplaint string             `json:"chief_complaint"`
	CreatedBy      EntityRef          `json:"created_by"`
	UpdatedBy      EntityRef          `json:"updated_by"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type VisitItemResponse struct {
	ID             uuid.UUID          `json:"id"`
	ExamineBy      EntityRef          `json:"examine_by"`
	Status         domain.VisitStatus `json:"status"`
	VisitDate      string             `json:"visit_date"`
	CheifComplaint string             `json:"chief_complaint"`
	CreatedBy      EntityRef          `json:"created_by"`
	UpdatedBy      EntityRef          `json:"updated_by"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type CreateVisitReq struct {
	PatientID      uuid.UUID          `json:"patient_id" binding:"required"`
	ExamineBy      uuid.UUID          `json:"examine_by" binding:"required"`
	Status         domain.VisitStatus `json:"status"`
	VisitDate      string             `json:"visit_date"`
	CheifComplaint string             `json:"chief_complaint"`
}

type PatientVisitsResponse struct {
	Patient PatientRef           `json:"patient"`
	Visits  []*VisitItemResponse `json:"visits"`
}

func VisitResponseFromDetails(v *domain.VisitDetails) *VisitResponse {
	return &VisitResponse{
		ID: v.ID,
		Patient: PatientRef{
			ID:   v.PatientID,
			Name: v.PatientName,
		},
		ExamineBy: EntityRef{
			ID:   v.ExamineBy,
			Name: v.ExamineByName,
		},
		Status:         v.Status,
		VisitDate:      v.VisitDate,
		CheifComplaint: v.CheifComplaint,
		CreatedBy: EntityRef{
			ID:   v.CreatedBy,
			Name: v.CreatedByName,
		},
		UpdatedBy: EntityRef{
			ID:   v.UpdatedBy,
			Name: v.UpdatedByName,
		},
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func visitItemFromDetails(v *domain.VisitDetails) *VisitItemResponse {
	return &VisitItemResponse{
		ID: v.ID,
		ExamineBy: EntityRef{
			ID:   v.ExamineBy,
			Name: v.ExamineByName,
		},
		Status:         v.Status,
		VisitDate:      v.VisitDate,
		CheifComplaint: v.CheifComplaint,
		CreatedBy: EntityRef{
			ID:   v.CreatedBy,
			Name: v.CreatedByName,
		},
		UpdatedBy: EntityRef{
			ID:   v.UpdatedBy,
			Name: v.UpdatedByName,
		},
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func PatientVisitsResponseFromDetails(vs []*domain.VisitDetails) *PatientVisitsResponse {
	if len(vs) == 0 {
		return &PatientVisitsResponse{
			Visits: []*VisitItemResponse{},
		}
	}

	visits := make([]*VisitItemResponse, 0, len(vs))

	for _, v := range vs {
		visits = append(visits, visitItemFromDetails(v))
	}

	return &PatientVisitsResponse{
		Patient: PatientRef{
			ID:   vs[0].PatientID,
			Name: vs[0].PatientName,
		},
		Visits: visits,
	}
}
