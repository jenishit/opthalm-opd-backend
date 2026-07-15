package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain/valueobjects"
)

type Role struct {
	ID       uuid.UUID
	RoleName string
}

type User struct {
	ID        uuid.UUID
	RoleID    uuid.UUID
	Email     string
	Password  valueobjects.Password
	CreatedAt time.Time
	UpdatedAt time.Time
}
type RoleUser struct {
	User
	Role Role
}

type UserRole string

type BasicDetails struct {
	ID       uuid.UUID             `db:"id"`
	UserRole UserRole              `db:"user_role"`
	Email    *string               `db:"email"`
	Username *string               `db:"username"`
	Password valueobjects.Password `db:"password"`
	FullName string                `db:"full_name"`
}
