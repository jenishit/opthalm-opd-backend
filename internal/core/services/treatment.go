package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type TreatmentService struct {
	repo port.TreatmentRepository
}

func NewTreatmentService(r port.TreatmentRepository) *TreatmentService {
	return &TreatmentService{repo: r}
}

func (s *TreatmentService) ReplaceTreatments(ctx context.Context, clinicID, visitID uuid.UUID, treatments []*domain.Treatment) ([]*domain.Treatment, error) {
	return s.repo.ReplaceTreatments(ctx, clinicID, visitID, treatments)
}

func (s *TreatmentService) ListTreatmentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.Treatment, error) {
	return s.repo.ListTreatmentsByVisitID(ctx, clinicID, visitID)
}
