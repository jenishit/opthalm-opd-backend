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

