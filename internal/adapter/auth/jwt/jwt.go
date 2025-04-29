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
