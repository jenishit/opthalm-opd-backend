package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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
			"email",
			"password",
		).
		Values(
			user.RoleID,
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
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, login *domain.Login) (*domain.BasicDetails, error) {
	var passwordHash string
	query, args, err := sq.
		Select(
			"u.id",
			"u.email",
			"u.password",
			"r.role_name",
		).
		From("users u").
		LeftJoin("role r on r.id = u.role_id").
		Where(sq.Eq{"u.email": login.Email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("User.FindByEmail build: %w", err)
	}

	var u domain.BasicDetails
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Email,
		&passwordHash,
		&u.UserRole,
	)
	password, err := valueobjects.NewPasswordFromHash(passwordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.CreateUser wrap password: %w", err)
	}
	u.Password = *password

	return &u, nil
}
