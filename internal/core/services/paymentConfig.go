package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type PaymentConfigService struct {
	repo port.PaymentConfigRepository
}

func NewPaymentConfigService(r port.PaymentConfigRepository) *PaymentConfigService {
	return &PaymentConfigService{repo: r}
}

func (s *PaymentConfigService) UpsertPaymentConfig(ctx context.Context, c *domain.PaymentConfig) (*domain.PaymentConfig, error) {
	return s.repo.UpsertPaymentConfig(ctx, c)
}

func (s *PaymentConfigService) GetPaymentConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PaymentConfig, error) {
	return s.repo.GetPaymentConfig(ctx, clinicID)
}

func (s *PaymentConfigService) SetQRAssetURL(ctx context.Context, clinicID uuid.UUID, url string) (*domain.PaymentConfig, error) {
	return s.repo.SetQRAssetURL(ctx, clinicID, url)
}
