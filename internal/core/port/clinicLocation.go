package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ClinicLocationRepository interface {
	CreateLocation(ctx context.Context, clinicID uuid.UUID, l *domain.ClinicLocation) (*domain.ClinicLocation, error)
	ListLocations(ctx context.Context, clinicID uuid.UUID) ([]*domain.ClinicLocation, error)
}

type ClinicLocationService interface {
	CreateLocation(ctx context.Context, clinicID uuid.UUID, l *domain.ClinicLocation) (*domain.ClinicLocation, error)
	ListLocations(ctx context.Context, clinicID uuid.UUID) ([]*domain.ClinicLocation, error)
}
