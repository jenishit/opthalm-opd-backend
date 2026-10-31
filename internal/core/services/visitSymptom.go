package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type VisitSymptomService struct {
	repo port.VisitSymptomRepository
}

func NewVisitSymptomService(r port.VisitSymptomRepository) *VisitSymptomService {
	return &VisitSymptomService{repo: r}
}

func (s *VisitSymptomService) UpsertVisitSymptom(ctx context.Context, clinicID uuid.UUID, vs *domain.VisitSymptom) (*domain.CreateVSResponse, error) {
	return s.repo.UpsertVisitSymptom(ctx, clinicID, vs)
}

func (s *VisitSymptomService) GetVisitSymptomByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitSymptom, error) {
	return s.repo.GetVisitSymptomByVisitID(ctx, clinicID, id)
}
