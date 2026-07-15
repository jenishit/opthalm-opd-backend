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
	// Valid bcrypt hashes start with $2a$, $2b$, or $2y$
	if len(hash) < 20 || (hash[0:3] != "$2a" && hash[0:3] != "$2b" && hash[0:3] != "$2y") {
		return nil, ErrInvalidHash
	}

	return &Password{
		hash: hash,
	}, nil
}

// ValidatePassword checks if a plaintext password meets security requirements.
// This is a public function for validation before creating a Password object.
func ValidatePassword(plaintext string) error {
	if plaintext == "" {
		return ErrPasswordEmpty
	}

	if len(plaintext) < PasswordMinLength {
		return ErrPasswordTooShort
	}

	return nil
}

// Hash returns the bcrypt hash of the password.
// This is used when storing the password in the database.
func (p *Password) Hash() string {
	return p.hash
}

// Verify checks if the provided plaintext password matches this Password's hash.
// Returns nil if the password matches, or an error if it doesn't.
func (p *Password) Verify(plaintext string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext)); err != nil {
		// Don't leak information about whether hash comparison failed or password was wrong
		return errors.New("password verification failed")
	}
	return nil
}

// Matches checks if another Password value object has the same hash.
// Two passwords are equal if their hashes are identical.
func (p *Password) Matches(other *Password) bool {
	if p == nil || other == nil {
		return p == other
	}
	return p.hash == other.hash
}

// String implements the Stringer interface but returns a redacted version for security.
// It never returns the actual hash to prevent accidental logging.
func (p *Password) String() string {
	if p == nil {
		return "<nil>"
	}
	return "Password(***)"
}

// MarshalJSON prevents the password hash from being accidentally serialized.
// This protects against accidental exposure in logs or API responses.
func (p *Password) MarshalJSON() ([]byte, error) {
	return []byte(`"Password(***)"` + "\n"), nil
}