package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type SubscriptionRepository interface {
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.Subscription, error)
	// Upsert creates or updates the one subscription row for a clinic
	// (clinic_id is unique) — used by the platform-operator endpoint to
	// activate/extend/cancel a clinic's subscription.
	Upsert(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
}

type SubscriptionService interface {
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.Subscription, error)
	Upsert(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
}
