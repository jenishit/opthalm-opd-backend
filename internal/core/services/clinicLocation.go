package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type ClinicLocationService struct {
	repo port.ClinicLocationRepository
}

func NewClinicLocationService(r port.ClinicLocationRepository) *ClinicLocationService {
	return &ClinicLocationService{repo: r}
}

func (s *ClinicLocationService) CreateLocation(ctx context.Context, clinicID uuid.UUID, l *domain.ClinicLocation) (*domain.ClinicLocation, error) {
	return s.repo.CreateLocation(ctx, clinicID, l)
}

func (s *ClinicLocationService) ListLocations(ctx context.Context, clinicID uuid.UUID) ([]*domain.ClinicLocation, error) {
	return s.repo.ListLocations(ctx, clinicID)
}
