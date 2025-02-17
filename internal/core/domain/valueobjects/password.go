package valueobjects

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordMinLength is the minimum required password length
const PasswordMinLength = 8

// PasswordHashCost is the bcrypt cost factor for hashing
const PasswordHashCost = 14

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordEmpty    = errors.New("password cannot be empty")
	ErrInvalidHash      = errors.New("invalid password hash")
)

// Password is a value object that encapsulates password logic following DDD principles.
// It is immutable and self-validating, ensuring password security and consistency.
type Password struct {
	// hash stores the bcrypt-hashed password (not the plaintext)
	hash string
}

// NewPassword creates a new Password value object by hashing the provided plaintext password.
// It validates the password meets minimum requirements before hashing.
//
