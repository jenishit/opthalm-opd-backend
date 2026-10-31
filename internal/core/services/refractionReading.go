package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type RefractionReadingService struct {
	repo port.RefractionReadingRepository
}

func NewRefractionReadingService(r port.RefractionReadingRepository) *RefractionReadingService {
	return &RefractionReadingService{repo: r}
}

func (s *RefractionReadingService) AddReading(ctx context.Context, clinicID uuid.UUID, r *domain.RefractionReading) (*domain.RefractionReading, error) {
	return s.repo.AddReading(ctx, clinicID, r)
}

func (s *RefractionReadingService) ListReadingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.RefractionReading, error) {
	return s.repo.ListReadingsByVisitID(ctx, clinicID, visitID)
}
