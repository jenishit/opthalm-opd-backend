package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain/valueobjects"
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
			"u.clinic_id",
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
		&u.ClinicID,
		&u.UserRole,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("UserRepo.GetUserByEmail scan: %w", err)
	}

	password, err := valueobjects.NewPasswordFromHash(passwordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.CreateUser wrap password: %w", err)
	}
	u.Password = *password

	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query, args, err := sq.
		Update("users").
		Set("password", passwordHash).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo.UpdatePassword build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("UserRepo.UpdatePassword exec: %w", err)
	}
	return nil
}

func (r *UserRepository) MarkEmailVerified(ctx context.Context, userID uuid.UUID) error {
	query, args, err := sq.
		Update("users").
		Set("email_verified_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo.MarkEmailVerified build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("UserRepo.MarkEmailVerified exec: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.BasicDetails, error) {
	var passwordHash string
	query, args, err := sq.
		Select(
			"u.id",
			"u.email",
			"u.password",
			"u.clinic_id",
			"r.role_name",
		).
		From("users u").
		LeftJoin("role r on r.id = u.role_id").
		Where(sq.Eq{"u.id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetUserByID build: %w", err)
	}

	var u domain.BasicDetails
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Email,
		&passwordHash,
		&u.ClinicID,
		&u.UserRole,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetUserByID scan: %w", err)
	}

	password, err := valueobjects.NewPasswordFromHash(passwordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetUserByID wrap password: %w", err)
	}
	u.Password = *password

	return &u, nil
}
