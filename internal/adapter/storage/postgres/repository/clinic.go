package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ClinicRepository struct {
	DB *postgres.DB
}

func NewClinicRepository(db *postgres.DB) *ClinicRepository {
	return &ClinicRepository{
		DB: db,
	}
}

func (sr *ClinicRepository) InsertClinic(ctx context.Context, s *domain.ClinicSettings) (*domain.ClinicSettings, error) {
	now := time.Now()

	query, args, err := sq.
		Insert("clinic_settings").
		Columns(
