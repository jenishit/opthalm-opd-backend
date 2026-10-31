package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type VisitDiagnosisRepository struct {
	DB *postgres.DB
}

func NewVisitDiagnosisRepository(db *postgres.DB) *VisitDiagnosisRepository {
	return &VisitDiagnosisRepository{DB: db}
}

// ReplaceVisitDiagnoses atomically swaps the full set of diagnoses recorded
// for a visit: every existing row is deleted and the given diagnoses are
// inserted in its place, so callers always PUT the complete set.
func (r *VisitDiagnosisRepository) ReplaceVisitDiagnoses(ctx context.Context, clinicID, visitID uuid.UUID, diagnoses []*domain.VisitDiagnosis) ([]*domain.VisitDiagnosis, error) {
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		delQuery, delArgs, err := sq.
			Delete("visit_diagnoses").
			Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("VisitDiagnosis.Replace build delete: %w", err)
		}

		if _, err := tx.Exec(ctx, delQuery, delArgs...); err != nil {
			return fmt.Errorf("VisitDiagnosis.Replace exec delete: %w", err)
		}

		for _, d := range diagnoses {
			query, args, err := sq.
				Insert("visit_diagnoses").
				Columns(
					"clinic_id",
					"visit_id",
					"diagnosis_id",
					"diagnosis_name",
					"icd10_code",
					"laterality",
					"status",
					"created_by",
				).
				Values(
					clinicID,
					visitID,
					d.DiagnosisID,
					d.DiagnosisName,
					d.Icd10Code,
					d.Laterality,
					d.Status,
					d.CreatedBy,
				).
				Suffix("RETURNING id, visit_id, created_at").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				return fmt.Errorf("VisitDiagnosis.Replace build insert: %w", err)
			}

			if err := tx.QueryRow(ctx, query, args...).Scan(&d.ID, &d.VisitID, &d.CreatedAt); err != nil {
				return fmt.Errorf("VisitDiagnosis.Replace exec insert: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return diagnoses, nil
}

func (r *VisitDiagnosisRepository) ListVisitDiagnosesByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.VisitDiagnosis, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"diagnosis_id",
			"diagnosis_name",
			"icd10_code",
			"laterality",
			"status",
			"created_by",
			"created_at",
		).
		From("visit_diagnoses").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		OrderBy("created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("VisitDiagnosis.ListByVisitID query: %w", err)
	}
	defer rows.Close()

	diagnoses := make([]*domain.VisitDiagnosis, 0)

	for rows.Next() {
		var d domain.VisitDiagnosis

		if err := rows.Scan(
			&d.ID,
			&d.VisitID,
			&d.DiagnosisID,
			&d.DiagnosisName,
			&d.Icd10Code,
			&d.Laterality,
			&d.Status,
			&d.CreatedBy,
			&d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("VisitDiagnosis.ListByVisitID scan: %w", err)
		}

		diagnoses = append(diagnoses, &d)
	}

	return diagnoses, nil
}
