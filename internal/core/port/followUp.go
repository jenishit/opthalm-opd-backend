package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type FollowUpRepository interface {
	UpsertFollowUp(ctx context.Context, clinicID uuid.UUID, f *domain.FollowUp) (*domain.FollowUp, error)
	GetFollowUpByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.FollowUp, error)
}

type FollowUpService interface {
	UpsertFollowUp(ctx context.Context, clinicID uuid.UUID, f *domain.FollowUp) (*domain.FollowUp, error)
	GetFollowUpByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.FollowUp, error)
}
