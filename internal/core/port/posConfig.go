package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PosConfigRepository interface {
	UpsertPosConfig(ctx context.Context, c *domain.PosConfig) (*domain.PosConfig, error)
	GetPosConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PosConfig, error)
}

type PosConfigService interface {
	UpsertPosConfig(ctx context.Context, c *domain.PosConfig) (*domain.PosConfig, error)
	GetPosConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PosConfig, error)
}
