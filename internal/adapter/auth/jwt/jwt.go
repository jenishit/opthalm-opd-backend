package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

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

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jt.secret), nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("access token expired")
			}
		}
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(string)
	roleName, ok2 := claims["role_name"].(string)
	if !ok || !ok2 {
		return nil, errors.New("invalid token claims")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user_id in token")
	}
	payload.UserId = uid
	payload.RoleName = roleName

	if sid, ok := claims["session_id"].(string); ok && sid != "" {
		sessID, err := uuid.Parse(sid)
		if err != nil {
			return nil, errors.New("invalid session_id in token")
		}
		payload.SessionID = sessID
	}

	return &payload, nil
}
