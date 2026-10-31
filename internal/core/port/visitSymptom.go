package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type VisitSymptomRepository interface {
	UpsertVisitSymptom(ctx context.Context, clinicID uuid.UUID, vs *domain.VisitSymptom) (*domain.CreateVSResponse, error)
	GetVisitSymptomByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitSymptom, error)
}

type VisitSymptomService interface {
	UpsertVisitSymptom(ctx context.Context, clinicID uuid.UUID, vs *domain.VisitSymptom) (*domain.CreateVSResponse, error)
	GetVisitSymptomByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitSymptom, error)
}
