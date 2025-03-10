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
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("MedicineRepo.Update build: %w", err)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MedicineRepo.Update exec: %w", err)
	}

	return nil
}

func (r *MedicineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Update("medicines").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("MedicineRepo.Delete build: %w", err)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MedicineRepo.Delete exec: %w", err)
	}

	return nil
}

// ─── Diagnosis Catalog ────────────────────────────────────────────────────────

type DiagnosisCatalogRepository struct {
	DB *postgres.DB
}

func NewDiagnosisCatalogRepository(db *postgres.DB) *DiagnosisCatalogRepository {
	return &DiagnosisCatalogRepository{DB: db}
}

func (r *DiagnosisCatalogRepository) Search(ctx context.Context, query string, limit int) ([]*domain.DiagnosisCatalog, error) {
	qb := sq.Select("id", "icd10_code", "name", "created_at", "updated_at").
		From("diagnosis_catalog").
		Where("deleted_at IS NULL").
		Where(sq.Or{
			sq.Expr("name ILIKE '%' || ? || '%'", query),
			sq.Expr("icd10_code ILIKE '%' || ? || '%'", query),
		}).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)

	return scanDiagnosisCatalogs(ctx, r.DB, qb)
}

func (r *DiagnosisCatalogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiagnosisCatalog, error) {
	qb := sq.Select("id", "icd10_code", "name", "created_at", "updated_at").
		From("diagnosis_catalog").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	return scanDiagnosisCatalog(ctx, r.DB, qb)
}

func (r *DiagnosisCatalogRepository) List(ctx context.Context, limit, offset int) ([]*domain.DiagnosisCatalog, error) {
	qb := sq.Select("id", "icd10_code", "name", "created_at", "updated_at").
		From("diagnosis_catalog").
		Where("deleted_at IS NULL").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	return scanDiagnosisCatalogs(ctx, r.DB, qb)
}

func (r *DiagnosisCatalogRepository) Update(ctx context.Context, d *domain.DiagnosisCatalog) error {
	query, args, err := sq.
		Update("diagnosis_catalog").
		Set("name", sq.Expr("COALESCE(?, name)", nullString(d.Name))).
		Set("icd10_code", sq.Expr("COALESCE(?, icd10_code)", nullStringPtr(d.Icd10Code))).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": d.ID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("DiagnosisCatalogRepo.Update build: %w", err)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DiagnosisCatalogRepo.Update exec: %w", err)
	}

	return nil
}

func (r *DiagnosisCatalogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Update("diagnosis_catalog").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("DiagnosisCatalogRepo.Delete build: %w", err)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DiagnosisCatalogRepo.Delete exec: %w", err)
	}

	return nil
}

// ─── History Condition ────────────────────────────────────────────────────────

type HistoryConditionRepository struct {
	DB *postgres.DB
}

func NewHistoryConditionRepository(db *postgres.DB) *HistoryConditionRepository {
	return &HistoryConditionRepository{DB: db}
}

func (r *HistoryConditionRepository) Search(ctx context.Context, query string, limit int) ([]*domain.HistoryCondition, error) {
	qb := sq.Select("id", "name", "created_at", "updated_at").
		From("history_conditions").
		Where("deleted_at IS NULL").
		Where(sq.Expr("name ILIKE '%' || ? || '%'", query)).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)

	return scanHistoryConditions(ctx, r.DB, qb)
}

func (r *HistoryConditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.HistoryCondition, error) {
	qb := sq.Select("id", "name", "created_at", "updated_at").
		From("history_conditions").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	return scanHistoryCondition(ctx, r.DB, qb)
}

func (r *HistoryConditionRepository) List(ctx context.Context, limit, offset int) ([]*domain.HistoryCondition, error) {
	qb := sq.Select("id", "name", "created_at", "updated_at").
		From("history_conditions").
		Where("deleted_at IS NULL").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	return scanHistoryConditions(ctx, r.DB, qb)
}

func (r *HistoryConditionRepository) Update(ctx context.Context, h *domain.HistoryCondition) error {
	query, args, err := sq.
		Update("history_conditions").
		Set("name", sq.Expr("COALESCE(?, name)", nullString(h.Name))).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": h.ID}).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("HistoryConditionRepo.Update build: %w", err)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("HistoryConditionRepo.Update exec: %w", err)
	}

	return nil
}

func (r *HistoryConditionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Update("history_conditions").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("HistoryConditionRepo.Delete build: %w", err)
	}

	_, err = r.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("HistoryConditionRepo.Delete exec: %w", err)
	}

	return nil
}

// ─── Scan helpers ─────────────────────────────────────────────────────────────

func scanMedicine(ctx context.Context, db *postgres.DB, qb sq.SelectBuilder) (*domain.Medicine, error) {
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var m domain.Medicine
	var brandName, strength, form sql.NullString

	err = db.QueryRow(ctx, query, args...).Scan(
		&m.ID,
		&m.MedicineName,
		&brandName,
		&strength,
		&form,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("scan: %w", err)
	}

	if brandName.Valid {
		m.BrandName = &brandName.String
	}
	if strength.Valid {
		m.Strength = &strength.String
	}
	if form.Valid {
		m.Form = &form.String
	}
