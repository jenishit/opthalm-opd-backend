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
			"gender",
			"occupation",
			"registered_on",
			"created_by",
			"updated_by",
		).
		Values(
			clinicID,
			pt.FullName,
			pt.Phone,
			pt.Address,
			pt.DOB,
			pt.Gender,
			pt.Occupation,
			sq.Expr("CURRENT_DATE"),
			pt.CreatedBy,
			pt.CreatedBy,
		).
		Suffix(`
	RETURNING
		id,
		full_name,
		registered_on,
		updated_by,
		created_at,
		updated_at
	`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("PatientRepo.CreatePatient build: %w", err)

	}

	err = pr.DB.QueryRow(ctx, query, args...).Scan(
		&pt.ID,
		&pt.FullName,
		&pt.RegisteredOn,
		&pt.UpdatedBy,
		&pt.CreatedAt,
		&pt.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("PatientRepo.CreatePatient scan: %w", err)
	}

	pt.ClinicID = clinicID

	return pt, nil
}

func (pr *PatientRepository) GetPatientByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.Patient, error) {
	var address, occupation sql.NullString

	query, args, err := sq.
		Select(
			"id",
			"full_name",
			"phone",
			"address",
			"dob",
			"gender",
			"occupation",
			"registered_on",
			"created_by",
			"created_at",
		).From("patients").
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
