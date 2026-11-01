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
	// SeatsLimit, PatientRecordsLimit, and LocationsLimit are the plan's
	// usage ceilings — set by the platform operator alongside PlanName, nil
	// meaning unlimited. The corresponding *Used counts are never stored
	// here; see SubscriptionUsage, which is computed live from the clinic's
	// actual rows so it can never drift out of sync.
	SeatsLimit          *int
	PatientRecordsLimit *int
	LocationsLimit      *int
	// PaymentMethodLabel is reserved for a future payment-processor
	// integration to sync a display label into (e.g. "Visa ending in
	// 4242") — never raw card data. Nothing writes to it yet.
	PaymentMethodLabel *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SubscriptionUsage is how many of each metered resource a clinic is
// currently using, computed live from its own data rather than stored —
// see Subscription's *Limit fields for the corresponding ceilings.
type SubscriptionUsage struct {
	SeatsUsed          int
	PatientRecordsUsed int
	LocationsUsed      int
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
