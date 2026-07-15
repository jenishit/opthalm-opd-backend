package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type RoleService struct {
	repo port.RoleRepository
}

func NewRoleService(rr port.RoleRepository) *RoleService {
	return &RoleService{
		repo: rr,
	}
}

func (r *RoleService) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	return r.repo.CreateRole(ctx, role)
}

func (r *RoleService) GetRoleIDByRoleName(ctx context.Context, role string) (*uuid.UUID, error) {
	return r.repo.GetRoleIDByRoleName(ctx, role)
}

func (r *RoleService) GetRoleNameByRoleID(ctx context.Context, id uuid.UUID) (*string, error) {
	return r.repo.GetRoleNameByRoleID(ctx, id)
}

