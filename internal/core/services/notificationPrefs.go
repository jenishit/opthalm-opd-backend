package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type NotificationPrefsService struct {
	repo port.NotificationPrefsRepository
}

func NewNotificationPrefsService(r port.NotificationPrefsRepository) *NotificationPrefsService {
	return &NotificationPrefsService{repo: r}
}

func (s *NotificationPrefsService) UpsertNotificationPrefs(ctx context.Context, p *domain.NotificationPrefs) (*domain.NotificationPrefs, error) {
	return s.repo.UpsertNotificationPrefs(ctx, p)
}

func (s *NotificationPrefsService) GetNotificationPrefs(ctx context.Context, clinicID uuid.UUID) (*domain.NotificationPrefs, error) {
	return s.repo.GetNotificationPrefs(ctx, clinicID)
}
