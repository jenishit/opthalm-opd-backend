package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error)
	GetRoleIDByRoleName(ctx context.Context, role string) (*uuid.UUID, error)
	GetRoleNameByRoleID(ctx context.Context, id uuid.UUID) (*string, error)
}

type RoleService interface {
	CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error)
	GetRoleIDByRoleName(ctx context.Context, role string) (*uuid.UUID, error)
	GetRoleNameByRoleID(ctx context.Context, id uuid.UUID) (*string, error)
}
