package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ExaminationFindingsRepository interface {
	UpsertExaminationFindings(ctx context.Context, clinicID uuid.UUID, ef *domain.ExaminationFindings) (*domain.ExaminationFindings, error)
	GetExaminationFindingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.ExaminationFindings, error)
}

type ExaminationFindingsService interface {
	UpsertExaminationFindings(ctx context.Context, clinicID uuid.UUID, ef *domain.ExaminationFindings) (*domain.ExaminationFindings, error)
	GetExaminationFindingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.ExaminationFindings, error)
}
