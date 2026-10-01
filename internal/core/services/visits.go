package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type VisitsService struct {
	repo port.VisitsRepository
}

func NewVisitsService(vr port.VisitsRepository) *VisitsService {
	return &VisitsService{
		repo: vr,
	}
}

func (vs *VisitsService) CreateVisit(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) (*domain.Visit, error) {
	return vs.repo.CreateVisit(ctx, clinicID, v)
}

func (vs *VisitsService) GetVisitByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitDetails, error) {
	return vs.repo.GetVisitByVisitID(ctx, clinicID, id)
}

func (vs *VisitsService) UpdateVisitByVisitID(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) error {
	return vs.repo.UpdateVisitByVisitID(ctx, clinicID, v)
}

func (vs *VisitsService) GetVisitsByPatientID(ctx context.Context, clinicID, id uuid.UUID) ([]*domain.VisitDetails, error) {
	return vs.repo.GetVisitsByPatientID(ctx, clinicID, id)
}
