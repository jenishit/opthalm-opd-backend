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
}

type SubscriptionResponse struct {
	ID               uuid.UUID `json:"id"`
	ClinicID         uuid.UUID `json:"clinic_id"`
	PlanName         string    `json:"plan_name"`
	Status           string    `json:"status"`
	CurrentPeriodEnd time.Time `json:"current_period_end"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func SubscriptionRes(s *domain.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:               s.ID,
		ClinicID:         s.ClinicID,
		PlanName:         s.PlanName,
		Status:           string(s.Status),
		CurrentPeriodEnd: s.CurrentPeriodEnd,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}
