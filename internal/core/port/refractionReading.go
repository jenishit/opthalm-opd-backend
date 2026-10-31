package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type RefractionReadingRepository interface {
	AddReading(ctx context.Context, clinicID uuid.UUID, r *domain.RefractionReading) (*domain.RefractionReading, error)
	ListReadingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.RefractionReading, error)
}

type RefractionReadingService interface {
	AddReading(ctx context.Context, clinicID uuid.UUID, r *domain.RefractionReading) (*domain.RefractionReading, error)
	ListReadingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.RefractionReading, error)
}
