package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VisitsRepository struct {
	DB *postgres.DB
}

func NewVisitsRepository(db *postgres.DB) *VisitsRepository {
	return &VisitsRepository{
		DB: db,
	}
}

func (vr *VisitsRepository) CreateVisit(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) (*domain.Visit, error) {
	status := v.Status
	if status == "" {
		status = domain.Scheduled
	}

	query, args, err := sq.
		Insert("visits").
		Columns(
			"clinic_id",
			"patient_id",
			"examined_by",
			"status",
			"chief_complaint",
			"created_by",
			"updated_by",
		).
		Values(
			clinicID,
			v.PatientID,
			v.ExamineBy,
			status,
			v.CheifComplaint,
			v.CreatedBy,
			v.UpdatedBy,
		).
		Suffix(`
			RETURNING
				id,
				status,
				created_at,
				updated_at
	`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("VisitRepo.CreateVisit build: %w", err)
	}

	err = vr.DB.QueryRow(ctx, query, args...).Scan(
		&v.ID,
		&v.Status,
		&v.CreatedAt,
		&v.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("VisitRepo.CreateVisit scan: %w", err)
	}

	v.ClinicID = clinicID

	return v, nil
}

func (vr *VisitsRepository) GetVisitByVisitID(ctx context.Context, clinicID, id uuid.UUID) (*domain.VisitDetails, error) {
	query, args, err := sq.
		Select(
			"v.id",
			"v.patient_id",
			"COALESCE(p.full_name, '') AS patient_name",
			"v.examined_by",
			"COALESCE(CONCAT(ep.first_name, ' ', ep.last_name), '') AS examined_by_name",
			"v.status",
			"v.visit_date",
			"v.chief_complaint",
			"v.created_by",
