package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, login *domain.Login) (*domain.BasicDetails, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.BasicDetails, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	MarkEmailVerified(ctx context.Context, userID uuid.UUID) error
}

type UserService interface {
	CreateUser(ctx context.Context, data *dto.CreateUser, clinicID uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, login *domain.Login) (*domain.BasicDetails, error)
}
