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

	accessToken, err := as.ts.CreateAccessToken(user, newSession.ID)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newPlain,
		SessionID:    newSession.ID,
		UserID:       user.ID,
		UserRole:     string(user.UserRole),
	}, nil
}

func (as *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := as.ts.HashRefreshToken(refreshToken)

	session, err := as.sessionRepo.GetByRefreshTokenHash(ctx, hash)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil
		}
		return err
	}

	return as.sessionRepo.Revoke(ctx, session.ID)
}

func (as *AuthService) issueTokenPair(ctx context.Context, user *domain.BasicDetails) (*domain.LoginResponse, error) {
	plain, hash, err := as.ts.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session, err := as.sessionRepo.Create(ctx, &domain.Session{
		UserID:           user.ID,
		RefreshTokenHash: hash,
		ExpiresAt:        time.Now().Add(as.refreshDuration),
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := as.ts.CreateAccessToken(user, session.ID)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: plain,
		SessionID:    session.ID,
		UserID:       user.ID,
		UserRole:     string(user.UserRole),
	}, nil
}

// RequestPasswordReset issues a password-reset token and emails it, if the
// address belongs to a user. It never reports whether the email exists —
// callers always get the same response — so this endpoint can't be used to
// enumerate registered accounts.
func (as *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := as.repo.GetUserByEmail(ctx, &domain.Login{Email: email})
	if err != nil {
		return nil
	}

	return as.issueVerificationToken(ctx, user.ID, domain.PurposePasswordReset, passwordResetTokenTTL, func(plain string) (string, string) {
		return "Reset your password",
			fmt.Sprintf("Use this code to reset your password: %s\n\nThis code expires in 1 hour. If you didn't request this, you can ignore this email.", plain)
	}, email)
}

// ConfirmPasswordReset applies a new password using a token minted by
// RequestPasswordReset, then revokes every existing session for that user —
// a password reset should force re-login everywhere, including on whatever
// device the attacker (if any) was using.
func (as *AuthService) ConfirmPasswordReset(ctx context.Context, plainToken, newPassword string) error {
	vt, err := as.consumeVerificationToken(ctx, plainToken, domain.PurposePasswordReset)
	if err != nil {
		return err
	}

	pwd, err := valueobjects.NewPassword(newPassword)
	if err != nil {
		return err
	}

	if err := as.repo.UpdatePassword(ctx, vt.UserID, pwd.Hash()); err != nil {
		return err
	}

	return as.sessionRepo.RevokeAllForUser(ctx, vt.UserID)
