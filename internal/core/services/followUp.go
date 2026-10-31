package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type FollowUpService struct {
	repo port.FollowUpRepository
}

func NewFollowUpService(r port.FollowUpRepository) *FollowUpService {
	return &FollowUpService{repo: r}
}

func (s *FollowUpService) UpsertFollowUp(ctx context.Context, clinicID uuid.UUID, f *domain.FollowUp) (*domain.FollowUp, error) {
	return s.repo.UpsertFollowUp(ctx, clinicID, f)
}

func (s *FollowUpService) GetFollowUpByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.FollowUp, error) {
	return s.repo.GetFollowUpByVisitID(ctx, clinicID, visitID)
}
