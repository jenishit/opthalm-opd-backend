package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PaymentConfigRepository interface {
	UpsertPaymentConfig(ctx context.Context, c *domain.PaymentConfig) (*domain.PaymentConfig, error)
	GetPaymentConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PaymentConfig, error)
	SetQRAssetURL(ctx context.Context, clinicID uuid.UUID, url string) (*domain.PaymentConfig, error)
}

type PaymentConfigService interface {
	UpsertPaymentConfig(ctx context.Context, c *domain.PaymentConfig) (*domain.PaymentConfig, error)
	GetPaymentConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PaymentConfig, error)
	SetQRAssetURL(ctx context.Context, clinicID uuid.UUID, url string) (*domain.PaymentConfig, error)
}
