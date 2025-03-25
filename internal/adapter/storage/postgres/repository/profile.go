package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ProfileRepository struct {
	DB *postgres.DB
}

func NewProfileRepository(db *postgres.DB) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (p *ProfileRepository) CreateProfile(ctx context.Context, pr *domain.Profile) (*domain.Profile, error) {
	query, args, err := sq.
		Insert("profile").
		Columns(
			"user_id",
			"first_name",
			"last_name",
			"phone",
