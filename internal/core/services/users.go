package services

import (
	"context"

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

func (u *UserService) CreateUser(ctx context.Context, data *dto.CreateUser) (*domain.User, error) {

	roleID, err := u.roleService.GetRoleIDByRoleName(ctx, data.RoleName)
	if err != nil {
		return nil, err
	}
	pwd, err := valueobjects.NewPassword(data.Password)
	if err != nil {
		return nil, err
	}
	us := &domain.User{
		Email:    data.Email,
		RoleID:   *roleID,
		Password: *pwd,
	}

	user, err := u.repo.CreateUser(ctx, us)
	if err != nil {
		return nil, err
	}

	profile := &domain.Profile{
		UserID:    user.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Phone:     &data.Phone,
	}
	_, err = u.profileService.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *UserService) GetUserByEmail(ctx context.Context, data *domain.Login) (*domain.BasicDetails, error) {
	return u.repo.GetUserByEmail(ctx, data)
}
