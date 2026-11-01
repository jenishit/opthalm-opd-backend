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

type SubscriptionRepository struct {
	DB *postgres.DB
}

func NewSubscriptionRepository(db *postgres.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.Subscription, error) {
	query, args, err := sq.Select(
		"id", "clinic_id", "plan_name", "status", "current_period_end",
		"seats_limit", "patient_records_limit", "locations_limit", "payment_method_label",
		"created_at", "updated_at",
	).
		From("subscriptions").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SubscriptionRepo.GetByClinicID build: %w", err)
	}

	var s domain.Subscription
	var status string
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&s.ID, &s.ClinicID, &s.PlanName, &status, &s.CurrentPeriodEnd,
		&s.SeatsLimit, &s.PatientRecordsLimit, &s.LocationsLimit, &s.PaymentMethodLabel,
		&s.CreatedAt, &s.UpdatedAt,
	)
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
		Columns(
			"clinic_id", "plan_name", "status", "current_period_end",
			"seats_limit", "patient_records_limit", "locations_limit",
		).
		Values(
			sub.ClinicID, sub.PlanName, string(sub.Status), sub.CurrentPeriodEnd,
			sub.SeatsLimit, sub.PatientRecordsLimit, sub.LocationsLimit,
		).
		Suffix(`
			ON CONFLICT (clinic_id) DO UPDATE SET
				plan_name = EXCLUDED.plan_name,
				status = EXCLUDED.status,
				current_period_end = EXCLUDED.current_period_end,
				seats_limit = EXCLUDED.seats_limit,
				patient_records_limit = EXCLUDED.patient_records_limit,
				locations_limit = EXCLUDED.locations_limit,
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

// GetUsageMetrics counts the clinic's current seats (users), patient
// records, and locations in one round trip.
func (r *SubscriptionRepository) GetUsageMetrics(ctx context.Context, clinicID uuid.UUID) (*domain.SubscriptionUsage, error) {
	var usage domain.SubscriptionUsage

	err := r.DB.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM users WHERE clinic_id = $1),
			(SELECT COUNT(*) FROM patients WHERE clinic_id = $1),
			(SELECT COUNT(*) FROM clinic_locations WHERE clinic_id = $1)
	`, clinicID).Scan(&usage.SeatsUsed, &usage.PatientRecordsUsed, &usage.LocationsUsed)
	if err != nil {
		return nil, fmt.Errorf("SubscriptionRepo.GetUsageMetrics scan: %w", err)
	}

	return &usage, nil
}
