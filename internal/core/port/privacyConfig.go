package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PrivacyConfigRepository interface {
	UpsertPrivacyConfig(ctx context.Context, c *domain.PrivacyConfig) (*domain.PrivacyConfig, error)
	GetPrivacyConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PrivacyConfig, error)
}

type PrivacyConfigService interface {
	UpsertPrivacyConfig(ctx context.Context, c *domain.PrivacyConfig) (*domain.PrivacyConfig, error)
	GetPrivacyConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PrivacyConfig, error)
}

type AuditLogRepository interface {
	RecordAuditLogEntry(ctx context.Context, e *domain.AuditLogEntry) (*domain.AuditLogEntry, error)
	ListAuditLog(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.AuditLogEntry, error)
}

type AuditLogService interface {
	RecordAuditLogEntry(ctx context.Context, e *domain.AuditLogEntry) (*domain.AuditLogEntry, error)
	ListAuditLog(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.AuditLogEntry, error)
}
