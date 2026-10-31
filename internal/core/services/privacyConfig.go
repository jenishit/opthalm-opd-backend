package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type PrivacyConfigService struct {
	repo port.PrivacyConfigRepository
}

func NewPrivacyConfigService(r port.PrivacyConfigRepository) *PrivacyConfigService {
	return &PrivacyConfigService{repo: r}
}

func (s *PrivacyConfigService) UpsertPrivacyConfig(ctx context.Context, c *domain.PrivacyConfig) (*domain.PrivacyConfig, error) {
	return s.repo.UpsertPrivacyConfig(ctx, c)
}

func (s *PrivacyConfigService) GetPrivacyConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PrivacyConfig, error) {
	return s.repo.GetPrivacyConfig(ctx, clinicID)
}

type AuditLogService struct {
	repo port.AuditLogRepository
}

func NewAuditLogService(r port.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: r}
}

func (s *AuditLogService) RecordAuditLogEntry(ctx context.Context, e *domain.AuditLogEntry) (*domain.AuditLogEntry, error) {
	return s.repo.RecordAuditLogEntry(ctx, e)
}

func (s *AuditLogService) ListAuditLog(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.AuditLogEntry, error) {
	return s.repo.ListAuditLog(ctx, clinicID, limit, offset)
}
