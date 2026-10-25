package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type PatientService struct {
	repo port.PatientRepository
}

func NewPatientService(pr port.PatientRepository) *PatientService {
	return &PatientService{
		repo: pr,
	}
}

func (ps *PatientService) CreatePatient(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) (*domain.Patient, error) {
	return ps.repo.CreatePatient(ctx, clinicID, pt)
}
func (ps *PatientService) GetPatientByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.Patient, error) {
	return ps.repo.GetPatientByID(ctx, clinicID, id)
}
func (ps *PatientService) GetPatients(ctx context.Context, clinicID uuid.UUID) ([]*domain.Patient, error) {
	return ps.repo.GetPatients(ctx, clinicID)
}

func (ps *PatientService) SearchPatients(ctx context.Context, clinicID uuid.UUID, query string, limit int) ([]*domain.Patient, error) {
	return ps.repo.SearchPatients(ctx, clinicID, query, limit)
}
func (ps *PatientService) UpdatePatientByID(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) error {
	return ps.repo.UpdatePatientByID(ctx, clinicID, pt)
}
func (ps *PatientService) DeletePatientByID(ctx context.Context, clinicID, id uuid.UUID) error {
	return ps.repo.DeletePatientByID(ctx, clinicID, id)
}
