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

func (ps *PatientService) CreatePatient(ctx context.Context, pt *domain.Patient) (*domain.Patient, error) {
	return ps.repo.CreatePatient(ctx, pt)
}
func (ps *PatientService) GetPatientByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	return ps.repo.GetPatientByID(ctx, id)
}
func (ps *PatientService) GetPatients(ctx context.Context) ([]*domain.Patient, error) {
	return ps.repo.GetPatients(ctx)
}
func (ps *PatientService) UpdatePatientByID(ctx context.Context, pt *domain.Patient) error {
	return ps.repo.UpdatePatientByID(ctx, pt)
}
func (ps *PatientService) DeletePatientByID(ctx context.Context, id uuid.UUID) error {
	return ps.repo.DeletePatientByID(ctx, id)
}
