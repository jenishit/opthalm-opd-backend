package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenPurpose distinguishes what a VerificationToken authorizes, since both
// password resets and email verification share the same table/shape.
type TokenPurpose string

const (
	PurposePasswordReset     TokenPurpose = "password_reset"
	PurposeEmailVerification TokenPurpose = "email_verification"
)

// VerificationToken is a single-use, expiring, high-entropy token issued for
// a specific purpose (password reset or email verification). Only its
// SHA-256 hash is ever persisted — the plaintext is emailed to the user once
// and never stored, the same pattern used for refresh tokens.
type VerificationToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Purpose   TokenPurpose
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ConfirmPasswordResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ConfirmEmailVerificationRequest struct {
	Token string `json:"token" binding:"required"`
}
