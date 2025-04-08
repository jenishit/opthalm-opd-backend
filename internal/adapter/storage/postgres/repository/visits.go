package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VisitsRepository struct {
	DB *postgres.DB
}

func NewVisitsRepository(db *postgres.DB) *VisitsRepository {
	return &VisitsRepository{
		DB: db,
	}
}

func (vr *VisitsRepository) CreateVisit(ctx context.Context, clinicID uuid.UUID, v *domain.Visit) (*domain.Visit, error) {
	status := v.Status
	if status == "" {
		status = domain.Scheduled
	}
