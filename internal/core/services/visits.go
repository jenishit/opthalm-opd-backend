package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type VisitsService struct {
	repo port.VisitsRepository
}

func NewVisitsService(vr port.VisitsRepository) *VisitsService {
	return &VisitsService{
		repo: vr,
	}
}

func (vs *VisitsService) CreateVisit(ctx context.Context, v *domain.Visit) (*domain.Visit, error) {
	return vs.repo.CreateVisit(ctx, v)
}

func (vs *VisitsService) GetVisitByVisitID(ctx context.Context, id uuid.UUID) (*domain.VisitDetails, error) {
	return vs.repo.GetVisitByVisitID(ctx, id)
}

func (vs *VisitsService) UpdateVisitByVisitID(ctx context.Context, v *domain.Visit) error {
	return vs.repo.UpdateVisitByVisitID(ctx, v)
}

func (vs *VisitsService) GetVisitsByPatientID(ctx context.Context, id uuid.UUID) ([]*domain.VisitDetails, error) {
	return vs.repo.GetVisitsByPatientID(ctx, id)
}
