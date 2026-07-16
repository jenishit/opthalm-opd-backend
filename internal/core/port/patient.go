package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type PatientRepository interface {
	CreatePatient(ctx context.Context, pt *domain.Patient) (*domain.Patient, error)
	GetPatientByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error)
	GetPatients(ctx context.Context) ([]*domain.Patient, error)
	UpdatePatientByID(ctx context.Context, pt *domain.Patient) error
	DeletePatientByID(ctx context.Context, pt *domain.Patient) error
}

type PatientService interface {
	CreatePatient(ctx context.Context, pt *domain.Patient) (*domain.Patient, error)
	GetPatientByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error)
	GetPatients(ctx context.Context) ([]*domain.Patient, error)
	UpdatePatientByID(ctx context.Context, pt *domain.Patient) error
	DeletePatientByID(ctx context.Context, pt *domain.Patient) error
}
