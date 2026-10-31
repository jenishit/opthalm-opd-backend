package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type NotificationPrefsRepository interface {
	UpsertNotificationPrefs(ctx context.Context, p *domain.NotificationPrefs) (*domain.NotificationPrefs, error)
	GetNotificationPrefs(ctx context.Context, clinicID uuid.UUID) (*domain.NotificationPrefs, error)
}

type NotificationPrefsService interface {
	UpsertNotificationPrefs(ctx context.Context, p *domain.NotificationPrefs) (*domain.NotificationPrefs, error)
	GetNotificationPrefs(ctx context.Context, clinicID uuid.UUID) (*domain.NotificationPrefs, error)
}
