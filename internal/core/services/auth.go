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
