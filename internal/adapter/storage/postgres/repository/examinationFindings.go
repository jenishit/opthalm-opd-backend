package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ExaminationFindingsRepository struct {
	DB *postgres.DB
}

func NewExaminationFindingsRepository(db *postgres.DB) *ExaminationFindingsRepository {
	return &ExaminationFindingsRepository{DB: db}
}

func (r *ExaminationFindingsRepository) UpsertExaminationFindings(ctx context.Context, clinicID uuid.UUID, ef *domain.ExaminationFindings) (*domain.ExaminationFindings, error) {
	query, args, err := sq.
		Insert("examination_findings").
		Columns(
			"clinic_id",
			"visit_id",
			"iop_od",
			"iop_os",
			"iop_time",
			"bp",
			"convergence",
			"cover_test",
			"eom",
			"anterior_segment_od",
			"anterior_segment_os",
			"fundus_od",
			"fundus_os",
			"created_by",
			"updated_by",
		).
		Values(
			clinicID,
			ef.VisitID,
			ef.IOPOD,
			ef.IOPOS,
			ef.IOPTime,
			ef.BP,
			ef.Convergence,
			ef.CoverTest,
			ef.EOM,
			ef.AnteriorSegmentOD,
			ef.AnteriorSegmentOS,
			ef.FundusOD,
			ef.FundusOS,
			ef.CreatedBy,
			ef.UpdatedBy,
		).
		Suffix(`ON CONFLICT (visit_id) DO UPDATE SET
			iop_od = EXCLUDED.iop_od,
			iop_os = EXCLUDED.iop_os,
			iop_time = EXCLUDED.iop_time,
			bp = EXCLUDED.bp,
			convergence = EXCLUDED.convergence,
			cover_test = EXCLUDED.cover_test,
			eom = EXCLUDED.eom,
			anterior_segment_od = EXCLUDED.anterior_segment_od,
			anterior_segment_os = EXCLUDED.anterior_segment_os,
			fundus_od = EXCLUDED.fundus_od,
			fundus_os = EXCLUDED.fundus_os,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
		RETURNING id, visit_id, created_by, updated_by, updated_at`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ExaminationFindings.Upsert build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&ef.ID,
		&ef.VisitID,
		&ef.CreatedBy,
		&ef.UpdatedBy,
		&ef.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ExaminationFindings.Upsert scan: %w", err)
	}

	return ef, nil
}

func (r *ExaminationFindingsRepository) GetExaminationFindingsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.ExaminationFindings, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"iop_od",
			"iop_os",
			"iop_time",
			"bp",
			"convergence",
			"cover_test",
			"eom",
			"anterior_segment_od",
			"anterior_segment_os",
			"fundus_od",
			"fundus_os",
			"created_by",
			"updated_by",
			"updated_at",
		).
		From("examination_findings").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var ef domain.ExaminationFindings

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&ef.ID,
		&ef.VisitID,
		&ef.IOPOD,
		&ef.IOPOS,
		&ef.IOPTime,
		&ef.BP,
		&ef.Convergence,
		&ef.CoverTest,
		&ef.EOM,
		&ef.AnteriorSegmentOD,
		&ef.AnteriorSegmentOS,
		&ef.FundusOD,
		&ef.FundusOS,
		&ef.CreatedBy,
		&ef.UpdatedBy,
		&ef.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("ExaminationFindings.GetByVisitID scan: %w", err)
	}

	return &ef, nil
}
