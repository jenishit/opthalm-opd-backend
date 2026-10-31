package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type RefractionReadingRepository struct {
	DB *postgres.DB
}

func NewRefractionReadingRepository(db *postgres.DB) *RefractionReadingRepository {
	return &RefractionReadingRepository{DB: db}
}

// AddReading upserts a reading keyed by (visit_id, eye, stage): re-measuring
// the same eye at the same stage replaces the earlier value instead of
// piling up duplicate rows.
func (r *RefractionReadingRepository) AddReading(ctx context.Context, clinicID uuid.UUID, rr *domain.RefractionReading) (*domain.RefractionReading, error) {
	query, args, err := sq.
		Insert("refraction_readings").
		Columns(
			"clinic_id",
			"visit_id",
			"eye",
			"stage",
			"sphere",
			"cylinder",
			"axis",
			"add_power",
			"va_achieved",
			"pd_mm",
			"lens_type",
			"created_by",
		).
		Values(
			clinicID,
			rr.VisitID,
			rr.Eye,
			rr.Stage,
			rr.Sphere,
			rr.Cylinder,
			rr.Axis,
			rr.AddPower,
			rr.VA,
			rr.PD,
			rr.LensType,
			rr.CreatedBy,
		).
		Suffix(`ON CONFLICT (visit_id, eye, stage) DO UPDATE SET
			sphere = EXCLUDED.sphere,
			cylinder = EXCLUDED.cylinder,
			axis = EXCLUDED.axis,
			add_power = EXCLUDED.add_power,
			va_achieved = EXCLUDED.va_achieved,
			pd_mm = EXCLUDED.pd_mm,
			lens_type = EXCLUDED.lens_type
		RETURNING id, visit_id, created_by, created_at`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("RefractionReading.AddReading build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&rr.ID,
		&rr.VisitID,
		&rr.CreatedBy,
		&rr.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("RefractionReading.AddReading scan: %w", err)
	}

	return rr, nil
}

func (r *RefractionReadingRepository) ListReadingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.RefractionReading, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"eye",
			"stage",
			"sphere",
			"cylinder",
			"axis",
			"add_power",
			"va_achieved",
			"pd_mm",
			"lens_type",
			"created_by",
			"created_at",
		).
		From("refraction_readings").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		OrderBy("created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("RefractionReading.ListByVisitID query: %w", err)
	}
	defer rows.Close()

	readings := make([]*domain.RefractionReading, 0)

	for rows.Next() {
		var rr domain.RefractionReading

		if err := rows.Scan(
			&rr.ID,
			&rr.VisitID,
			&rr.Eye,
			&rr.Stage,
			&rr.Sphere,
			&rr.Cylinder,
			&rr.Axis,
			&rr.AddPower,
			&rr.VA,
			&rr.PD,
			&rr.LensType,
			&rr.CreatedBy,
			&rr.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("RefractionReading.ListByVisitID scan: %w", err)
		}

		readings = append(readings, &rr)
	}

	return readings, nil
}
