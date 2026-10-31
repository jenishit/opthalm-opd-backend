package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type TokenService interface {
	CreateAccessToken(user *domain.BasicDetails, sessionID uuid.UUID) (string, error)
	VerifyAccessToken(token string) (*domain.TokenPayload, error)
	// GenerateRefreshToken returns a high-entropy plaintext refresh token and
	// its SHA-256 hash (the hash is what gets stored; the plaintext is
	// returned to the client once and never persisted).
	GenerateRefreshToken() (plain string, hash string, err error)
	// HashRefreshToken hashes a presented refresh token the same way, for
	// looking it up against a stored hash.
	HashRefreshToken(plain string) string
}

type AuthService interface {
	Login(ctx context.Context, details *domain.Login) (*domain.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Signup(ctx context.Context, req *domain.SignupRequest) (*domain.LoginResponse, error)
	RequestPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, plainToken, newPassword string) error
	RequestEmailVerification(ctx context.Context, userID uuid.UUID) error
	ConfirmEmailVerification(ctx context.Context, plainToken string) error
}
