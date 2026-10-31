package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type PatientHistoryRepository interface {
	CreatePatientHistory(ctx context.Context, clinicID uuid.UUID, h *domain.PatientHistory) (*domain.PatientHistory, error)
	ListPatientHistoryByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.PatientHistory, error)
}

type PatientHistoryService interface {
	CreatePatientHistory(ctx context.Context, clinicID uuid.UUID, h *domain.PatientHistory) (*domain.PatientHistory, error)
	ListPatientHistoryByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.PatientHistory, error)
}
