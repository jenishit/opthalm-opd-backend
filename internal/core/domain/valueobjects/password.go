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
// Returns an error if:
//   - password is empty
//   - password is shorter than PasswordMinLength (8 chars)
//   - bcrypt hashing fails
func NewPassword(plaintext string) (*Password, error) {
	if err := ValidatePassword(plaintext); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), PasswordHashCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &Password{
		hash: string(hash),
	}, nil
}

// NewPasswordFromHash creates a Password value object from an existing bcrypt hash.
// This is useful when loading passwords from the database.
//
// Returns an error if:
//   - hash is empty
//   - hash is not a valid bcrypt hash
func NewPasswordFromHash(hash string) (*Password, error) {
	if hash == "" {
		return nil, ErrInvalidHash
	}

	// Validate it's a valid bcrypt hash by attempting to use it
