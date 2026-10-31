package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type VisitDiagnosisService struct {
	repo port.VisitDiagnosisRepository
}

func NewVisitDiagnosisService(r port.VisitDiagnosisRepository) *VisitDiagnosisService {
	return &VisitDiagnosisService{repo: r}
}

func (s *VisitDiagnosisService) ReplaceVisitDiagnoses(ctx context.Context, clinicID, visitID uuid.UUID, diagnoses []*domain.VisitDiagnosis) ([]*domain.VisitDiagnosis, error) {
	return s.repo.ReplaceVisitDiagnoses(ctx, clinicID, visitID, diagnoses)
}

func (s *VisitDiagnosisService) ListVisitDiagnosesByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.VisitDiagnosis, error) {
	return s.repo.ListVisitDiagnosesByVisitID(ctx, clinicID, visitID)
}
