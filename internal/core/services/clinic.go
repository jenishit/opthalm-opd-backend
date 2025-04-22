package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type ClinicService struct {
	repo port.ClinicRepository
}

func NewClinicService(cr port.ClinicRepository) *ClinicService {
	return &ClinicService{
		repo: cr,
	}
}

func (cs *ClinicService) InsertClinic(ctx context.Context, s *domain.ClinicSettings) (*domain.ClinicSettings, error) {
	return cs.repo.InsertClinic(ctx, s)
}

func (cs *ClinicService) GetClinicByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicSettings, error) {
	return cs.repo.GetClinicByClinicID(ctx, clinicID)
}

func (cs *ClinicService) UpdateClinic(ctx context.Context, s *domain.ClinicSettings) error {
	return cs.repo.UpdateClinic(ctx, s)
}

func (cs *ClinicService) GetAllClinics(ctx context.Context) ([]*domain.ClinicSettings, error) {
	return cs.repo.GetAllClinics(ctx)
}
