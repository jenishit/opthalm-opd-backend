package repository

import (
	"context"
	"errors"
	"fmt"

	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type FollowUpRepository struct {
	DB *postgres.DB
}

func NewFollowUpRepository(db *postgres.DB) *FollowUpRepository {
	return &FollowUpRepository{DB: db}
}

func (r *FollowUpRepository) UpsertFollowUp(ctx context.Context, clinicID uuid.UUID, f *domain.FollowUp) (*domain.FollowUp, error) {
	query, args, err := sq.
		Insert("follow_ups").
		Columns(
			"clinic_id",
			"visit_id",
			"followup_date",
			"reason",
			"created_by",
			"updated_by",
		).
		Values(
			clinicID,
			f.VisitID,
			f.FollowUpDate,
			f.Reason,
			f.CreatedBy,
			f.UpdatedBy,
		).
		Suffix(`ON CONFLICT (visit_id) DO UPDATE SET
			followup_date = EXCLUDED.followup_date,
			reason = EXCLUDED.reason,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
		RETURNING id, visit_id, reminder_sent_at, created_by, updated_by, updated_at`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("FollowUp.Upsert build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&f.ID,
		&f.VisitID,
		&f.ReminderSentAt,
		&f.CreatedBy,
		&f.UpdatedBy,
		&f.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("FollowUp.Upsert scan: %w", err)
	}

	return f, nil
}

func (r *FollowUpRepository) GetFollowUpByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) (*domain.FollowUp, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"followup_date",
			"reason",
			"reminder_sent_at",
			"created_by",
			"updated_by",
			"updated_at",
		).
		From("follow_ups").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var f domain.FollowUp
	var followUpDate time.Time

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&f.ID,
		&f.VisitID,
		&followUpDate,
		&f.Reason,
		&f.ReminderSentAt,
		&f.CreatedBy,
		&f.UpdatedBy,
		&f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("FollowUp.GetByVisitID scan: %w", err)
	}

	f.FollowUpDate = followUpDate.Format("2006-01-02")

	return &f, nil
}
