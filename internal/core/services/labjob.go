package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type LabJobService struct {
	repo port.LabJobRepository
}

func NewLabJobService(r port.LabJobRepository) *LabJobService {
	return &LabJobService{repo: r}
}

func (s *LabJobService) Create(ctx context.Context, clinicID uuid.UUID, job *domain.LabJob) (*domain.LabJobDetails, error) {
	return s.repo.Create(ctx, clinicID, job)
}

func (s *LabJobService) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.LabJobDetails, error) {
	return s.repo.GetByID(ctx, clinicID, id)
}

func (s *LabJobService) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.LabJobDetails, error) {
	return s.repo.List(ctx, clinicID, limit, offset)
}

func (s *LabJobService) ListByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.LabJobDetails, error) {
	return s.repo.ListByPatientID(ctx, clinicID, patientID)
}

func (s *LabJobService) UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.LabJobStatus, notes *string, changedBy uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, clinicID, id, status, notes, changedBy)
}
