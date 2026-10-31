package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type PatientHistoryService struct {
	repo port.PatientHistoryRepository
}

func NewPatientHistoryService(r port.PatientHistoryRepository) *PatientHistoryService {
	return &PatientHistoryService{repo: r}
}

func (s *PatientHistoryService) CreatePatientHistory(ctx context.Context, clinicID uuid.UUID, h *domain.PatientHistory) (*domain.PatientHistory, error) {
	return s.repo.CreatePatientHistory(ctx, clinicID, h)
}

func (s *PatientHistoryService) ListPatientHistoryByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.PatientHistory, error) {
	return s.repo.ListPatientHistoryByPatientID(ctx, clinicID, patientID)
}
