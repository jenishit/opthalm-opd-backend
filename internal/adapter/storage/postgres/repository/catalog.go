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
		PlaceholderFormat(sq.Dollar)

	return scanMedicines(ctx, r.DB, qb)
}

func (r *MedicineRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Medicine, error) {
	qb := sq.Select("id", "medicine_name", "brand_name", "strength", "form", "created_at", "updated_at").
		From("medicines").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	return scanMedicine(ctx, r.DB, qb)
}

func (r *MedicineRepository) List(ctx context.Context, limit, offset int) ([]*domain.Medicine, error) {
	qb := sq.Select("id", "medicine_name", "brand_name", "strength", "form", "created_at", "updated_at").
		From("medicines").
		Where("deleted_at IS NULL").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	return scanMedicines(ctx, r.DB, qb)
}

func (r *MedicineRepository) Update(ctx context.Context, m *domain.Medicine) error {
	builder := sq.Update("medicines").
		Set("medicine_name", sq.Expr("COALESCE(?, medicine_name)", nullString(m.MedicineName))).
		Set("brand_name", sq.Expr("COALESCE(?, brand_name)", nullStringPtr(m.BrandName))).
		Set("strength", sq.Expr("COALESCE(?, strength)", nullStringPtr(m.Strength))).
		Set("form", sq.Expr("COALESCE(?, form)", nullStringPtr(m.Form))).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": m.ID}).
