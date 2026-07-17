package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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

func (pr *PatientRepository) CreatePatient(ctx context.Context, pt *domain.Patient) (*domain.Patient, error) {
	query, args, err := sq.
		Insert("patients").
		Columns(
			"full_name",
			"phone",
			"address",
			"dob",
			"gender",
			"occupation",
			"created_by",
		).
		Values(
			pt.FullName,
			pt.Phone,
			pt.Address,
			pt.DOB,
			pt.Gender,
			pt.Occupation,
			pt.CreatedBy,
		).
		Suffix(`
	RETURNING
		id,
		full_name,
		created_at
	`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("PatientRepo.CreatePatient build: %w", err)

	}

	err = pr.DB.QueryRow(ctx, query, args...).Scan(
		&pt.ID,
		&pt.FullName,
		&pt.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("PatientRepo.CreatePatient scan: %w", err)
	}

	return pt, nil
}

func (pr *PatientRepository) GetPatientByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
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
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var patient domain.Patient

	err = pr.DB.QueryRow(ctx, query, args...).Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Phone,
		&address,
		&patient.DOB,
		&patient.Gender,
		&occupation,
		&patient.RegisteredOn,
		&patient.CreatedBy,
		&patient.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	if address.Valid {
		patient.Address = &address.String
	}
	if occupation.Valid {
		patient.Occupation = &occupation.String
	}

	return &patient, nil
}

func (pr *PatientRepository) GetPatients(ctx context.Context) ([]*domain.Patient, error) {
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
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := pr.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var patients []*domain.Patient

	for rows.Next() {
		var p domain.Patient

		err := rows.Scan(
			&p.ID,
			&p.FullName,
			&p.Phone,
			&address,
			&p.DOB,
			&p.Gender,
			&occupation,
			&p.RegisteredOn,
			&p.CreatedBy,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if address.Valid {
			p.Address = &address.String
		}
		if occupation.Valid {
			p.Occupation = &occupation.String
		}

		patients = append(patients, &p)
	}

	return patients, nil
}

func (pr *PatientRepository) UpdatePatientByID(ctx context.Context, pt *domain.Patient) error {
	name := nullString(pt.FullName)
	phone := nullString(pt.Phone)
	address := nullString(*pt.Address)
	dob := nullString(pt.DOB)
	gender := nullString(pt.Gender)
	occupation := nullString(*pt.Occupation)

	query, args, err := sq.
		Update("patients").
		Set("full_name", sq.Expr("COALESCE(?, full_name)", name)).
		Set("phone", sq.Expr("COALESCE(?, phone)", phone)).
		Set("address", sq.Expr("COALESCE(?, address)", address)).
		Set("dob", sq.Expr("COALESCE(?, dob)", dob)).
		Set("gender", sq.Expr("COALESCE(?, gender)", gender)).
		Set("occupation", sq.Expr("COALESCE(?, occupation)", occupation)).
		Set("updated_at", sq.Expr("NOW()")).
		Set("updated_by", sq.Expr("updated_by", pt.UpdatedBy)).
		Where(sq.Eq{"id": pt.ID}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = pr.DB.Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("Failed to update result: %w", err)
	}

	return nil
}

func (pr *PatientRepository) DeletePatientByID(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Update("patients").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = pr.DB.Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("Failed to Delete patient: %w", err)
	}

	return nil
}
