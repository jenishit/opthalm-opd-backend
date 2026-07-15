package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenPayload struct {
	UserId    uuid.UUID `json:"user_id"`
	RoleName  string    `json:"role_name"`
	SessionID uuid.UUID `json:"session_id"`
}

type RefreshTokenPayload struct {
	UserId    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"role_name"`
}
type SessionState struct {
	SessionID           uuid.UUID `json:"session_id"`
	UserID              uuid.UUID `json:"user_id"`
	RoleID              uuid.UUID `json:"role_id"`
	Email               string    `json:"email"`
	RoleName            string    `json:"role_name"`
	IssuedAt            time.Time `json:"issued_at"`
	ExpiresAt           time.Time `json:"expires_at"`
}
