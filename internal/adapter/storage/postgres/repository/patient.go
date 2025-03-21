package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type PatientRepository struct {
	DB *postgres.DB
}

func NewPatientRepository(db *postgres.DB) *PatientRepository {
	return &PatientRepository{
		DB: db,
	}
}

func (pr *PatientRepository) CreatePatient(ctx context.Context, clinicID uuid.UUID, pt *domain.Patient) (*domain.Patient, error) {
	query, args, err := sq.
		Insert("patients").
		Columns(
			"clinic_id",
			"full_name",
			"phone",
			"address",
			"dob",
