package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) (*domain.Session, error)
	GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error)
	// Rotate atomically revokes the old session (setting replaced_by) and
	// inserts the new one.
	Rotate(ctx context.Context, oldSessionID uuid.UUID, newSession *domain.Session) (*domain.Session, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
