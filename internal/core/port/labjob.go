package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type LabJobRepository interface {
	Create(ctx context.Context, clinicID uuid.UUID, job *domain.LabJob) (*domain.LabJobDetails, error)
	GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.LabJobDetails, error)
	List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.LabJobDetails, error)
	ListByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.LabJobDetails, error)
	// UpdateStatus updates the job's status and inserts a status_history row
	// in the same transaction.
	UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.LabJobStatus, notes *string, changedBy uuid.UUID) error
}

type LabJobService interface {
	Create(ctx context.Context, clinicID uuid.UUID, job *domain.LabJob) (*domain.LabJobDetails, error)
	GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.LabJobDetails, error)
	List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.LabJobDetails, error)
	ListByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.LabJobDetails, error)
	UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.LabJobStatus, notes *string, changedBy uuid.UUID) error
}
