package port

import (
	"context"

	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, login *domain.Login) (*domain.BasicDetails, error)
}

type UserService interface {
	CreateUser(ctx context.Context, data *dto.CreateUser) (*domain.User, error)
	GetUserByEmail(ctx context.Context, login *domain.Login) (*domain.BasicDetails, error)
}
