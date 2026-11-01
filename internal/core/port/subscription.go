package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type SubscriptionRepository interface {
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.Subscription, error)
	// Upsert creates or updates the one subscription row for a clinic
	// (clinic_id is unique) — used by the platform-operator endpoint to
	// activate/extend/cancel a clinic's subscription.
	Upsert(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	// GetUsageMetrics counts the clinic's current seats/patient
	// records/locations live, for comparison against the subscription's
	// *Limit fields — never stored, always fresh.
	GetUsageMetrics(ctx context.Context, clinicID uuid.UUID) (*domain.SubscriptionUsage, error)
}

type SubscriptionService interface {
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.Subscription, error)
	Upsert(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	GetUsageMetrics(ctx context.Context, clinicID uuid.UUID) (*domain.SubscriptionUsage, error)
}
