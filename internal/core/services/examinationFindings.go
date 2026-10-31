package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type ExaminationFindingsService struct {
	repo port.ExaminationFindingsRepository
}

func NewExaminationFindingsService(r port.ExaminationFindingsRepository) *ExaminationFindingsService {
	return &ExaminationFindingsService{repo: r}
}

func (s *ExaminationFindingsService) UpsertExaminationFindings(ctx context.Context, clinicID uuid.UUID, ef *domain.ExaminationFindings) (*domain.ExaminationFindings, error) {
	return s.repo.UpsertExaminationFindings(ctx, clinicID, ef)
}

func (s *ExaminationFindingsService) GetExaminationFindingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.ExaminationFindings, error) {
	return s.repo.GetExaminationFindingsByVisitID(ctx, clinicID, visitID)
}
