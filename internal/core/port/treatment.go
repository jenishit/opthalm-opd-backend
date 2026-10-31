package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type TreatmentRepository interface {
	ReplaceTreatments(ctx context.Context, clinicID, visitID uuid.UUID, treatments []*domain.Treatment) ([]*domain.Treatment, error)
	ListTreatmentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.Treatment, error)
}

type TreatmentService interface {
	ReplaceTreatments(ctx context.Context, clinicID, visitID uuid.UUID, treatments []*domain.Treatment) ([]*domain.Treatment, error)
	ListTreatmentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.Treatment, error)
}
