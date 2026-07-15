package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain/valueobjects"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

func NewAuthService(userRepo port.UserRepository, tokenService port.TokenService) *AuthService {
	return &AuthService{
		repo: userRepo,
		ts:   tokenService,
	}
}

func (as *AuthService) Login(ctx context.Context, details *domain.Login) (*domain.LoginResponse, error) {
	sessionID := uuid.New()
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

	accessToken, err := as.ts.CreateAccessToken(user, sessionID)

	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken: accessToken,
		SessionID:   sessionID,
		UserID:      user.ID,
		UserRole:    string(user.UserRole),
	}, nil
}
