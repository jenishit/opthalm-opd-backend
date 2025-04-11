package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain/valueobjects"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

// passwordResetTokenTTL/emailVerificationTokenTTL bound how long an issued
// reset/verification link stays usable.
const (
	passwordResetTokenTTL     = 1 * time.Hour
	emailVerificationTokenTTL = 24 * time.Hour
)

type AuthService struct {
	repo             port.UserRepository
	sessionRepo      port.SessionRepository
	signupRepo       port.SignupRepository
	verificationRepo port.VerificationTokenRepository
	emailSender      port.EmailSender
	ts               port.TokenService
	refreshDuration  time.Duration
}

