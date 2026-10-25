package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VerificationTokenRepository interface {
	Create(ctx context.Context, t *domain.VerificationToken) (*domain.VerificationToken, error)
	GetByHash(ctx context.Context, hash string) (*domain.VerificationToken, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	// InvalidateAllForUser marks every not-yet-used token of the given
	// purpose for this user as used, so requesting a new reset/verification
	// link invalidates any older, still-pending one.
	InvalidateAllForUser(ctx context.Context, userID uuid.UUID, purpose domain.TokenPurpose) error
}

// EmailSender abstracts email delivery so AuthService doesn't depend on SMTP
// directly. The default wiring (main.go) falls back to a log-only sender
// when no SMTP server is configured, so password-reset/verification flows
// work end-to-end in dev without any mail setup.
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}
