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

type PatientHistoryRepository struct {
	DB *postgres.DB
}

func NewPatientHistoryRepository(db *postgres.DB) *PatientHistoryRepository {
	return &PatientHistoryRepository{DB: db}
}

// CreatePatientHistory appends a new history entry and links it to its
// catalog conditions. History is a log, so this never updates an existing
// row — every call inserts a fresh one.
func (r *PatientHistoryRepository) CreatePatientHistory(ctx context.Context, clinicID uuid.UUID, h *domain.PatientHistory) (*domain.PatientHistory, error) {
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		query, args, err := sq.
			Insert("patient_history").
			Columns(
				"clinic_id",
				"patient_id",
				"uses_spectacles",
				"systemic_medications",
				"topical_medications",
				"notes",
				"created_by",
			).
			Values(
				clinicID,
				h.PatientID,
				h.UsesSpectacles,
				h.SystemicMedications,
				h.TopicalMedications,
				h.Notes,
				h.CreatedBy,
			).
			Suffix("RETURNING id, patient_id, created_by, created_at").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("PatientHistory.Create build insert: %w", err)
		}

		if err := tx.QueryRow(ctx, query, args...).Scan(&h.ID, &h.PatientID, &h.CreatedBy, &h.CreatedAt); err != nil {
			return fmt.Errorf("PatientHistory.Create exec insert: %w", err)
		}

		for _, conditionID := range h.ConditionIDs {
			query, args, err := sq.
				Insert("patient_history_conditions").
				Columns("history_id", "condition_id").
				Values(h.ID, conditionID).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				return fmt.Errorf("PatientHistory.Create build condition link: %w", err)
			}

			if _, err := tx.Exec(ctx, query, args...); err != nil {
				return fmt.Errorf("PatientHistory.Create exec condition link: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}

// ListPatientHistoryByPatientID returns every history entry for a patient,
// newest first, with each entry's condition_ids attached.
func (r *PatientHistoryRepository) ListPatientHistoryByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.PatientHistory, error) {
	query, args, err := sq.
		Select(
			"id",
			"patient_id",
			"uses_spectacles",
			"systemic_medications",
			"topical_medications",
			"notes",
			"created_by",
			"created_at",
		).
		From("patient_history").
		Where(sq.Eq{"patient_id": patientID, "clinic_id": clinicID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PatientHistory.ListByPatientID query: %w", err)
	}
	defer rows.Close()

	history := make([]*domain.PatientHistory, 0)
	byID := make(map[uuid.UUID]*domain.PatientHistory)
	ids := make([]uuid.UUID, 0)

	for rows.Next() {
		var h domain.PatientHistory

		if err := rows.Scan(
			&h.ID,
			&h.PatientID,
			&h.UsesSpectacles,
			&h.SystemicMedications,
			&h.TopicalMedications,
			&h.Notes,
			&h.CreatedBy,
			&h.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("PatientHistory.ListByPatientID scan: %w", err)
		}
		h.ConditionIDs = []uuid.UUID{}

		history = append(history, &h)
		byID[h.ID] = &h
		ids = append(ids, h.ID)
	}

	if len(ids) == 0 {
		return history, nil
	}

	conditionQuery, conditionArgs, err := sq.
		Select("history_id", "condition_id").
		From("patient_history_conditions").
		Where(sq.Eq{"history_id": ids}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	conditionRows, err := r.DB.Query(ctx, conditionQuery, conditionArgs...)
	if err != nil {
		return nil, fmt.Errorf("PatientHistory.ListByPatientID query conditions: %w", err)
	}
	defer conditionRows.Close()

	for conditionRows.Next() {
		var historyID, conditionID uuid.UUID
		if err := conditionRows.Scan(&historyID, &conditionID); err != nil {
			return nil, fmt.Errorf("PatientHistory.ListByPatientID scan condition: %w", err)
		}

		if h, ok := byID[historyID]; ok {
			h.ConditionIDs = append(h.ConditionIDs, conditionID)
		}
	}

	return history, nil
}
