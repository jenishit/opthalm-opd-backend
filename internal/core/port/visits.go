package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type VisitsRepository interface {
	CreateVisit(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) (*domain.Visit, error)
	GetVisitByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitDetails, error)
	UpdateVisitByVisitID(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) error
	GetVisitsByPatientID(ctx context.Context, clinicID, id uuid.UUID) ([]*domain.VisitDetails, error)
}

type VisitsService interface {
	CreateVisit(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) (*domain.Visit, error)
	GetVisitByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitDetails, error)
	UpdateVisitByVisitID(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) error
	GetVisitsByPatientID(ctx context.Context, clinicID, id uuid.UUID) ([]*domain.VisitDetails, error)
}
