package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type IntegrationService struct {
	repo port.IntegrationRepository
}

func NewIntegrationService(r port.IntegrationRepository) *IntegrationService {
	return &IntegrationService{repo: r}
}

func (s *IntegrationService) ListIntegrations(ctx context.Context, clinicID uuid.UUID) ([]*domain.Integration, error) {
	return s.repo.ListIntegrations(ctx, clinicID)
}

func (s *IntegrationService) ConnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) (*domain.Integration, error) {
	return s.repo.ConnectIntegration(ctx, clinicID, provider)
}

func (s *IntegrationService) DisconnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) error {
	return s.repo.DisconnectIntegration(ctx, clinicID, provider)
}
