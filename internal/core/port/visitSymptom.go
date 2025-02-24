package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VisitSymptomRepository interface {
	UpsertVisitSymptom(ctx context.Context, vs *domain.VisitSymptom) (*domain.CreateVSResponse, error)
	GetVisitSymptomByVisitID(ctx context.Context, id uuid.UUID) (*domain.VisitSymptom, error)
}

type VisitSymptomService interface {
	UpsertVisitSymptom(ctx context.Context, vs *domain.VisitSymptom) (*domain.CreateVSResponse, error)
	GetVisitSymptomByVisitID(ctx context.Context, id uuid.UUID) (*domain.VisitSymptom, error)
}
