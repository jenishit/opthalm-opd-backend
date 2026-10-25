package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionTrialing SubscriptionStatus = "trialing"
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionPastDue  SubscriptionStatus = "past_due"
	SubscriptionCanceled SubscriptionStatus = "cancelled"
)

type Subscription struct {
	ID               uuid.UUID
	ClinicID         uuid.UUID
	PlanName         string
	Status           SubscriptionStatus
	CurrentPeriodEnd time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// IsActive reports whether this subscription currently grants access —
// trialing or active, and not past its current period end.
func (s *Subscription) IsActive() bool {
	if s == nil {
		return false
	}
	if s.Status != SubscriptionTrialing && s.Status != SubscriptionActive {
		return false
	}
	return time.Now().Before(s.CurrentPeriodEnd)
}
