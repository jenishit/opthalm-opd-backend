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
)

type RoleRepository struct {
	DB *postgres.DB
}

func NewRoleRepository(db *postgres.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	query, args, err := sq.
		Insert("role").
		Columns("role_name").
		Values(role.RoleName).
		Suffix(`RETURNING id`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("RoleRepo.CreateRole build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&role.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("RoleRepo.CreateRole scan: %w", err)
	}

	return role, nil
}

func (r *RoleRepository) GetRoleIDByRoleName(ctx context.Context, role string) (*uuid.UUID, error) {
	query, args, err := sq.
		Select("id").
		From("role").
		Where(sq.Eq{"role_name": role}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("RoleRepo.GetRoleIDByRoleName build: %w", err)
	}
	var roleID uuid.UUID
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&roleID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
