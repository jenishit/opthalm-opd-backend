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

type VisitCreateResponse struct {
	ID             uuid.UUID          `json:"id"`
	PatientID      uuid.UUID          `json:"patient_id"`
	ExamineBy      uuid.UUID          `json:"examine_by"`
	Status         domain.VisitStatus `json:"status"`
	CheifComplaint string             `json:"chief_complaint"`
	CreatedBy      uuid.UUID          `json:"created_by"`
	UpdatedBy      uuid.UUID          `json:"updated_by"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
