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

type TreatmentRepository struct {
	DB *postgres.DB
}

func NewTreatmentRepository(db *postgres.DB) *TreatmentRepository {
	return &TreatmentRepository{DB: db}
}

// ReplaceTreatments atomically swaps the full set of treatments recorded for
// a visit: every existing row is deleted and the given treatments are
// inserted in its place, so callers always PUT the complete set.
func (r *TreatmentRepository) ReplaceTreatments(ctx context.Context, clinicID, visitID uuid.UUID, treatments []*domain.Treatment) ([]*domain.Treatment, error) {
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		delQuery, delArgs, err := sq.
			Delete("treatments").
			Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("Treatment.Replace build delete: %w", err)
		}

		if _, err := tx.Exec(ctx, delQuery, delArgs...); err != nil {
			return fmt.Errorf("Treatment.Replace exec delete: %w", err)
		}

		for _, t := range treatments {
			query, args, err := sq.
				Insert("treatments").
				Columns(
					"clinic_id",
					"visit_id",
					"type",
					"medicine_id",
					"medicine_name",
					"procedure_name",
					"advice_text",
					"route",
					"dosage",
					"frequency",
					"duration",
					"created_by",
				).
				Values(
					clinicID,
					visitID,
					t.Type,
					t.MedicineID,
					t.MedicineName,
					t.ProcedureName,
					t.AdviceText,
					t.Route,
					t.Dosage,
					t.Frequency,
					t.Duration,
					t.CreatedBy,
				).
				Suffix("RETURNING id, visit_id, created_at").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				return fmt.Errorf("Treatment.Replace build insert: %w", err)
			}

			if err := tx.QueryRow(ctx, query, args...).Scan(&t.ID, &t.VisitID, &t.CreatedAt); err != nil {
				return fmt.Errorf("Treatment.Replace exec insert: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return treatments, nil
}

func (r *TreatmentRepository) ListTreatmentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.Treatment, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"type",
			"medicine_id",
			"medicine_name",
			"procedure_name",
			"advice_text",
			"route",
			"dosage",
			"frequency",
			"duration",
			"created_by",
			"created_at",
		).
		From("treatments").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		OrderBy("created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Treatment.ListByVisitID query: %w", err)
	}
	defer rows.Close()

	treatments := make([]*domain.Treatment, 0)

	for rows.Next() {
		var t domain.Treatment

		if err := rows.Scan(
			&t.ID,
			&t.VisitID,
			&t.Type,
			&t.MedicineID,
			&t.MedicineName,
			&t.ProcedureName,
			&t.AdviceText,
			&t.Route,
			&t.Dosage,
			&t.Frequency,
			&t.Duration,
			&t.CreatedBy,
			&t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("Treatment.ListByVisitID scan: %w", err)
		}

		treatments = append(treatments, &t)
	}

	return treatments, nil
}
