package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ClinicRepository interface {
	InsertClinic(ctx context.Context, s *domain.ClinicSettings) (*domain.ClinicSettings, error)
	GetClinicByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicSettings, error)
	UpdateClinic(ctx context.Context, s *domain.ClinicSettings) error
	GetAllClinics(ctx context.Context) ([]*domain.ClinicSettings, error)
}

type ClinicService interface {
	InsertClinic(ctx context.Context, s *domain.ClinicSettings) (*domain.ClinicSettings, error)
	GetClinicByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicSettings, error)
	UpdateClinic(ctx context.Context, s *domain.ClinicSettings) error
	GetAllClinics(ctx context.Context) ([]*domain.ClinicSettings, error)
}
 