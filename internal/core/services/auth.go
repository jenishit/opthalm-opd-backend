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

func NewAuthService(
	userRepo port.UserRepository,
	sessionRepo port.SessionRepository,
	signupRepo port.SignupRepository,
	verificationRepo port.VerificationTokenRepository,
	emailSender port.EmailSender,
	tokenService port.TokenService,
	refreshDuration time.Duration,
) *AuthService {
	return &AuthService{
		repo:             userRepo,
		sessionRepo:      sessionRepo,
		signupRepo:       signupRepo,
		verificationRepo: verificationRepo,
		emailSender:      emailSender,
		ts:               tokenService,
		refreshDuration:  refreshDuration,
	}
}

func (as *AuthService) Signup(ctx context.Context, req *domain.SignupRequest) (*domain.LoginResponse, error) {
	pwd, err := valueobjects.NewPassword(req.AdminPassword)
	if err != nil {
		return nil, err
	}

	user, err := as.signupRepo.CreateTenant(ctx, req, pwd.Hash())
	if err != nil {
		return nil, err
	}

	return as.issueTokenPair(ctx, user)
}

func (as *AuthService) Login(ctx context.Context, details *domain.Login) (*domain.LoginResponse, error) {
	user, err := as.repo.GetUserByEmail(ctx, details)
	if err != nil {
		return nil, err
	}

	passwordVO, err := valueobjects.NewPasswordFromHash(user.Password.Hash())
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if err := passwordVO.Verify(details.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return as.issueTokenPair(ctx, user)
}

func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
	hash := as.ts.HashRefreshToken(refreshToken)

	session, err := as.sessionRepo.GetByRefreshTokenHash(ctx, hash)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	if session.RevokedAt != nil {
		// This refresh token was already rotated away and is being presented
		// again — a strong signal it was stolen. Revoke every session for
		// this user so both the thief and the legitimate holder are forced
		// to log in again.
		_ = as.sessionRepo.RevokeAllForUser(ctx, session.UserID)
		return nil, domain.ErrInvalidToken
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrExpiredToken
	}

	user, err := as.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	newPlain, newHash, err := as.ts.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newSession := &domain.Session{
		UserID:           user.ID,
		RefreshTokenHash: newHash,
		ExpiresAt:        time.Now().Add(as.refreshDuration),
	}
	newSession, err = as.sessionRepo.Rotate(ctx, session.ID, newSession)
	if err != nil {
		return nil, err
	}
