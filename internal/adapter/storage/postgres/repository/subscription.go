package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type SubscriptionRepository struct {
	DB *postgres.DB
}

func NewSubscriptionRepository(db *postgres.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.Subscription, error) {
	query, args, err := sq.Select("id", "clinic_id", "plan_name", "status", "current_period_end", "created_at", "updated_at").
		From("subscriptions").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SubscriptionRepo.GetByClinicID build: %w", err)
	}

	var s domain.Subscription
	var status string
	err = r.DB.QueryRow(ctx, query, args...).Scan(&s.ID, &s.ClinicID, &s.PlanName, &status, &s.CurrentPeriodEnd, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("SubscriptionRepo.GetByClinicID scan: %w", err)
	}
	s.Status = domain.SubscriptionStatus(status)

	return &s, nil
}

func (r *SubscriptionRepository) Upsert(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	query, args, err := sq.Insert("subscriptions").
		Columns("clinic_id", "plan_name", "status", "current_period_end").
		Values(sub.ClinicID, sub.PlanName, string(sub.Status), sub.CurrentPeriodEnd).
		Suffix(`
			ON CONFLICT (clinic_id) DO UPDATE SET
				plan_name = EXCLUDED.plan_name,
				status = EXCLUDED.status,
				current_period_end = EXCLUDED.current_period_end,
				updated_at = NOW()
			RETURNING id, created_at, updated_at
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SubscriptionRepo.Upsert build: %w", err)
	}

	if err := r.DB.QueryRow(ctx, query, args...).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
		return nil, fmt.Errorf("SubscriptionRepo.Upsert exec: %w", err)
	}

	return sub, nil
}
