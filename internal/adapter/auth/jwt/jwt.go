package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

// refreshTokenBytes is the entropy of a generated refresh token (32 bytes =
// 256 bits, well beyond what's brute-forceable).
const refreshTokenBytes = 32

type JWTToken struct {
	secret   string
	duration time.Duration
}

func New(config *config.Token) (port.TokenService, error) {
	durationStr := config.Duration
	duration, err := time.ParseDuration(durationStr)

	if err != nil {
		return nil, err
	}

	if config.Secret == "" {
		return nil, errors.New("JWT secret is missing")
	}

	return &JWTToken{
		secret:   config.Secret,
		duration: duration,
	}, nil
}

func (jt *JWTToken) CreateAccessToken(user *domain.BasicDetails, sessionID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(jt.duration)

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"role_name":  user.UserRole,
		"clinic_id":  user.ClinicID,
		"exp":        expirationTime.Unix(),
		"session_id": sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jt.secret))
	if err != nil {
		return "", fmt.Errorf("creating access token: %w", err)
	}

	return tokenString, nil
}

func (jt *JWTToken) VerifyAccessToken(tokenString string) (*domain.TokenPayload, error) {
	var payload domain.TokenPayload
