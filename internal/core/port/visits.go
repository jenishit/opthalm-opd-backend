package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VisitsRepository interface {
	CreateVisit(ctx context.Context, v *domain.Visit) (*domain.Visit, error)
	GetVisitByVisitID(ctx context.Context, id uuid.UUID) (*domain.VisitDetails, error)
	UpdateVisitByVisitID(ctx context.Context, v *domain.Visit) error
	GetVisitsByPatientID(ctx context.Context, id uuid.UUID) ([]*domain.VisitDetails, error)
}

type VisitsService interface {
	CreateVisit(ctx context.Context, v *domain.Visit) (*domain.Visit, error)
	GetVisitByVisitID(ctx context.Context, id uuid.UUID) (*domain.VisitDetails, error)
	UpdateVisitByVisitID(ctx context.Context, v *domain.Visit) error
	GetVisitsByPatientID(ctx context.Context, id uuid.UUID) ([]*domain.VisitDetails, error)
}
