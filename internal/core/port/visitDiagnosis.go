package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VisitDiagnosisRepository interface {
	ReplaceVisitDiagnoses(ctx context.Context, clinicID, visitID uuid.UUID, diagnoses []*domain.VisitDiagnosis) ([]*domain.VisitDiagnosis, error)
	ListVisitDiagnosesByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.VisitDiagnosis, error)
}

type VisitDiagnosisService interface {
	ReplaceVisitDiagnoses(ctx context.Context, clinicID, visitID uuid.UUID, diagnoses []*domain.VisitDiagnosis) ([]*domain.VisitDiagnosis, error)
	ListVisitDiagnosesByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.VisitDiagnosis, error)
}
