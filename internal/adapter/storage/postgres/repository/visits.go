package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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

func (vr *VisitsRepository) CreateVisit(ctx context.Context, v *domain.Visit) (*domain.Visit, error) {
	query, args, err := sq.
		Insert("visits").
		Columns(
			"patient_id",
			"examined_by",
			"status",
			"chief_complaint",
		).
		Values(
			v.PatientID,
			v.ExamineBy,
			v.Status,
			v.CheifComplaint,
		).
		Suffix(` 
			RETURNING
				id,
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
		&v.CreatedAt,
		&v.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ProfileRepo.CreateProfile scan: %w", err)
	}

	return v, nil
}

func (vr *VisitsRepository) GetVisitByVisitID(ctx context.Context, id uuid.UUID) (*domain.VisitDetails, error) {
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
		Where(sq.Eq{"v.id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var v domain.VisitDetails

	err = vr.DB.QueryRow(ctx, query, args...).Scan(
		&v.ID,
		&v.PatientID,
		&v.PatientName,
		&v.ExamineBy,
		&v.ExamineByName,
		&v.Status,
		&v.VisitDate,
		&v.CheifComplaint,
		&v.CreatedBy,
		&v.CreatedByName,
		&v.UpdatedBy,
		&v.UpdatedByName,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &v, nil
}

func (vr *VisitsRepository) UpdateVisitByVisitID(ctx context.Context, v *domain.Visit) error {
	query, args, err := sq.
		Update("visits").
		Set("examined_by", sq.Expr("COALESCE(?, examined_by)", v.ExamineBy)).
		Set("status", sq.Expr("COALESCE(?, status)", v.Status)).
		Set("visit_date", sq.Expr("COALESCE(?, visit_date)", v.VisitDate)).
		Set("chief_complaint", sq.Expr("COALESCE(?, chief_complaint)", v.CheifComplaint)).
		Set("updated_at", sq.Expr("NOW()")).
		Set("updated_by", v.UpdatedBy).
		Where(sq.Eq{"id": v.ID}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("VisitRepo.UpdateVisitByVisitID build: %w", err)
	}

	_, err = vr.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("VisitRepo.UpdateVisitByVisitID exec: %w", err)
	}

	return nil
}

func (vr *VisitsRepository) GetVisitsByPatientID(ctx context.Context, id uuid.UUID) ([]*domain.VisitDetails, error) {
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
		Where(sq.Eq{"v.patient_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := vr.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var visits []*domain.VisitDetails

	for rows.Next() {
		var v domain.VisitDetails

		err := rows.Scan(
			&v.ID,
			&v.PatientID,
			&v.PatientName,
			&v.ExamineBy,
			&v.ExamineByName,
			&v.Status,
			&v.VisitDate,
			&v.CheifComplaint,
			&v.CreatedBy,
			&v.CreatedByName,
			&v.UpdatedBy,
			&v.UpdatedByName,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		visits = append(visits, &v)
	}

	return visits, nil
}
