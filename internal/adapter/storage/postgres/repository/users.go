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
			"email",
			"password",
		).
		Values(
			user.RoleID,
			user.ClinicID,
			user.Email,
			user.Password.Hash(),
		).
		//The returning data when a user is created is placed in the suffix
		Suffix(`
			RETURNING
				id,
				created_at,
				updated_at
		`). //To be safe from SQLInjection, the variables are replaced with dollar
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("UserRepo.CreateUser build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
