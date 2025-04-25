package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain/valueobjects"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type UserService struct {
	repo           port.UserRepository
	roleService    port.RoleService
	profileService port.ProfileService
}

func NewUserService(ur port.UserRepository, rsvc port.RoleService, psvc port.ProfileService) *UserService {
	return &UserService{
		repo:           ur,
		roleService:    rsvc,
		profileService: psvc,
	}
}

func (u *UserService) CreateUser(ctx context.Context, data *dto.CreateUser, clinicID uuid.UUID) (*domain.User, error) {

	roleID, err := u.roleService.GetRoleIDByRoleName(ctx, data.RoleName)
	if err != nil {
		return nil, err
	}
	pwd, err := valueobjects.NewPassword(data.Password)
