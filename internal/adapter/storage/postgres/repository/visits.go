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
			"COALESCE(CONCAT(cp.first_name, ' ', cp.last_name), '') AS created_by_name",
			"v.updated_by",
			"COALESCE(CONCAT(up.first_name, ' ', up.last_name), '') AS updated_by_name",
			"v.created_at",
			"v.updated_at",
		).
		From("visits v").
		LeftJoin("patients p ON p.id = v.patient_id").
		LeftJoin("users eu ON eu.id = v.examined_by").
		LeftJoin("profile ep ON ep.user_id = eu.id").
		LeftJoin("users cu ON cu.id = v.created_by").
		LeftJoin("profile cp ON cp.user_id = cu.id").
		LeftJoin("users uu ON uu.id = v.updated_by").
		LeftJoin("profile up ON up.user_id = uu.id").
		Where(sq.Eq{"v.id": id, "v.clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var v domain.VisitDetails
	var visitDate time.Time

	err = vr.DB.QueryRow(ctx, query, args...).Scan(
		&v.ID,
		&v.PatientID,
		&v.PatientName,
		&v.ExamineBy,
		&v.ExamineByName,
		&v.Status,
		&visitDate,
		&v.CheifComplaint,
		&v.CreatedBy,
		&v.CreatedByName,
		&v.UpdatedBy,
		&v.UpdatedByName,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	v.VisitDate = visitDate.Format(time.RFC3339)

	return &v, nil
}

func (vr *VisitsRepository) UpdateVisitByVisitID(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) error {
	query, args, err := sq.
		Update("visits").
		Set("examined_by", sq.Expr("COALESCE(?, examined_by)", nullUUID(v.ExamineBy))).
		Set("status", sq.Expr("COALESCE(?, status)", nullString(string(v.Status)))).
		Set("visit_date", sq.Expr("COALESCE(?, visit_date)", nullString(v.VisitDate))).
