package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	UserAgent        *string
	IPAddress        *string
	IssuedAt         time.Time
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	ReplacedBy       *uuid.UUID
}
