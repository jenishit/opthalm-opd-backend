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

// ─── Medicines ────────────────────────────────────────────────────────────────

type MedicineRepository struct {
	DB *postgres.DB
}

func NewMedicineRepository(db *postgres.DB) *MedicineRepository {
	return &MedicineRepository{DB: db}
}

func (r *MedicineRepository) Search(ctx context.Context, query string, limit int) ([]*domain.Medicine, error) {
	qb := sq.Select("id", "medicine_name", "brand_name", "strength", "form", "created_at", "updated_at").
		From("medicines").
		Where("deleted_at IS NULL").
		Where(sq.Or{
			sq.Expr("medicine_name ILIKE '%' || ? || '%'", query),
			sq.Expr("brand_name ILIKE '%' || ? || '%'", query),
		}).
		Limit(uint64(limit)).
