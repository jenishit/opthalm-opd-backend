package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type PosConfigService struct {
	repo port.PosConfigRepository
}

func NewPosConfigService(r port.PosConfigRepository) *PosConfigService {
	return &PosConfigService{repo: r}
}

func (s *PosConfigService) UpsertPosConfig(ctx context.Context, c *domain.PosConfig) (*domain.PosConfig, error) {
	return s.repo.UpsertPosConfig(ctx, c)
}

func (s *PosConfigService) GetPosConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PosConfig, error) {
	return s.repo.GetPosConfig(ctx, clinicID)
}
