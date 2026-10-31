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
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/config"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
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
				return nil, domain.ErrExpiredToken
			}
		}
		return nil, domain.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	roleName, ok2 := claims["role_name"].(string)
	if !ok || !ok2 {
		return nil, domain.ErrInvalidToken
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	payload.UserId = uid
	payload.RoleName = roleName

	if cid, ok := claims["clinic_id"].(string); ok && cid != "" {
		clinicID, err := uuid.Parse(cid)
		if err != nil {
			return nil, domain.ErrInvalidToken
		}
		payload.ClinicID = clinicID
	}

	if sid, ok := claims["session_id"].(string); ok && sid != "" {
		sessID, err := uuid.Parse(sid)
		if err != nil {
			return nil, domain.ErrInvalidToken
		}
		payload.SessionID = sessID
	}

	return &payload, nil
}

func (jt *JWTToken) GenerateRefreshToken() (plain string, hash string, err error) {
	buf := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}
	plain = base64.RawURLEncoding.EncodeToString(buf)
	return plain, jt.HashRefreshToken(plain), nil
}

func (jt *JWTToken) HashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
