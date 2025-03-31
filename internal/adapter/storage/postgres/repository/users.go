package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain/valueobjects"
)

type UserRepository struct {
	DB *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query, args, err := sq.
		Insert("users").
		Columns(
			"role_id",
			"clinic_id",
