package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type IntegrationRepository interface {
	ListIntegrations(ctx context.Context, clinicID uuid.UUID) ([]*domain.Integration, error)
	ConnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) (*domain.Integration, error)
	DisconnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) error
}

type IntegrationService interface {
	ListIntegrations(ctx context.Context, clinicID uuid.UUID) ([]*domain.Integration, error)
	ConnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) (*domain.Integration, error)
	DisconnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) error
}
