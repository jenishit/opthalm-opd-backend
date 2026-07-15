package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type TokenService interface {
	CreateAccessToken(user *domain.BasicDetails, sessionID uuid.UUID) (string, error)
	VerifyAccessToken(token string) (*domain.TokenPayload, error)
}

type AuthService interface {
	Login(ctx context.Context, details *domain.Login) (*domain.LoginResponse, error)

}