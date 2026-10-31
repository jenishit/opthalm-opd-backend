package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type IntegrationResponse struct {
	ID          uuid.UUID  `json:"id"`
	ClinicID    uuid.UUID  `json:"clinic_id"`
	Provider    string     `json:"provider"`
	Connected   bool       `json:"connected"`
	ConnectedAt *time.Time `json:"connected_at,omitempty"`
}

func IntegrationRes(i *domain.Integration) *IntegrationResponse {
	if i == nil {
		return nil
	}

	return &IntegrationResponse{
		ID:          i.ID,
		ClinicID:    i.ClinicID,
		Provider:    i.Provider,
		Connected:   i.Connected,
		ConnectedAt: i.ConnectedAt,
	}
}

func IntegrationResList(integrations []*domain.Integration) []*IntegrationResponse {
	rsp := make([]*IntegrationResponse, 0, len(integrations))
	for _, i := range integrations {
		rsp = append(rsp, IntegrationRes(i))
	}
	return rsp
}
