package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type UpsertSubscriptionReq struct {
	PlanName         string `json:"plan_name"`
	Status           string `json:"status" binding:"required,oneof=trialing active past_due cancelled"`
	CurrentPeriodEnd string `json:"current_period_end" binding:"required"` // YYYY-MM-DD
	// SeatsLimit, PatientRecordsLimit, and LocationsLimit are the plan's
	// usage ceilings; omit (or send null) for unlimited.
	SeatsLimit          *int `json:"seats_limit" binding:"omitempty,min=0"`
	PatientRecordsLimit *int `json:"patient_records_limit" binding:"omitempty,min=0"`
	LocationsLimit      *int `json:"locations_limit" binding:"omitempty,min=0"`
}

type SubscriptionResponse struct {
	ID                  uuid.UUID `json:"id"`
	ClinicID            uuid.UUID `json:"clinic_id"`
	PlanName            string    `json:"plan_name"`
	Status              string    `json:"status"`
	CurrentPeriodEnd    time.Time `json:"current_period_end"`
	SeatsLimit          *int      `json:"seats_limit,omitempty"`
	PatientRecordsLimit *int      `json:"patient_records_limit,omitempty"`
	LocationsLimit      *int      `json:"locations_limit,omitempty"`
	// PaymentMethodLabel is reserved for a future payment-processor
	// integration — never raw card data, and nothing sets it yet.
	PaymentMethodLabel *string   `json:"payment_method_label,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func SubscriptionRes(s *domain.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:                  s.ID,
		ClinicID:            s.ClinicID,
		PlanName:            s.PlanName,
		Status:              string(s.Status),
		CurrentPeriodEnd:    s.CurrentPeriodEnd,
		SeatsLimit:          s.SeatsLimit,
		PatientRecordsLimit: s.PatientRecordsLimit,
		LocationsLimit:      s.LocationsLimit,
		PaymentMethodLabel:  s.PaymentMethodLabel,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}

// BillingSubscriptionResponse is what a clinic's own admin sees on the
// Billing & subscription tab: the subscription plus how much of each
// metered resource it's currently using, computed live.
type BillingSubscriptionResponse struct {
	SubscriptionResponse
	SeatsUsed          int `json:"seats_used"`
	PatientRecordsUsed int `json:"patient_records_used"`
	LocationsUsed      int `json:"locations_used"`
}

func BillingSubscriptionRes(s *domain.Subscription, usage *domain.SubscriptionUsage) *BillingSubscriptionResponse {
	return &BillingSubscriptionResponse{
		SubscriptionResponse: *SubscriptionRes(s),
		SeatsUsed:            usage.SeatsUsed,
		PatientRecordsUsed:   usage.PatientRecordsUsed,
		LocationsUsed:        usage.LocationsUsed,
	}
}
