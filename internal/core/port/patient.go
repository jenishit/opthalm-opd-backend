package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PatientRepository interface {
	CreatePatient(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) (*domain.Patient, error)
	GetPatientByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.Patient, error)
	GetPatients(ctx context.Context, clinicID uuid.UUID) ([]*domain.Patient, error)
	SearchPatients(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.Patient, error)
	UpdatePatientByID(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) error
	DeletePatientByID(ctx context.Context, clinicID, id uuid.UUID) error
}

type PatientService interface {
	CreatePatient(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) (*domain.Patient, error)
	GetPatientByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.Patient, error)
	GetPatients(ctx context.Context, clinicID uuid.UUID) ([]*domain.Patient, error)
	SearchPatients(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.Patient, error)
	UpdatePatientByID(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) error
	DeletePatientByID(ctx context.Context, clinicID, id uuid.UUID) error
}
