package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type VerificationTokenRepository struct {
	DB *postgres.DB
}

func NewVerificationTokenRepository(db *postgres.DB) *VerificationTokenRepository {
	return &VerificationTokenRepository{DB: db}
}

func (r *VerificationTokenRepository) Create(ctx context.Context, t *domain.VerificationToken) (*domain.VerificationToken, error) {
	query, args, err := sq.Insert("verification_tokens").
		Columns("user_id", "purpose", "token_hash", "expires_at").
		Values(t.UserID, string(t.Purpose), t.TokenHash, t.ExpiresAt).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("VerificationTokenRepo.Create build: %w", err)
	}
	if err := r.DB.QueryRow(ctx, query, args...).Scan(&t.ID, &t.CreatedAt); err != nil {
		return nil, fmt.Errorf("VerificationTokenRepo.Create exec: %w", err)
	}
	return t, nil
}

func (r *VerificationTokenRepository) GetByHash(ctx context.Context, hash string) (*domain.VerificationToken, error) {
	query, args, err := sq.Select("id", "user_id", "purpose", "token_hash", "expires_at", "used_at", "created_at").
		From("verification_tokens").
		Where(sq.Eq{"token_hash": hash}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("VerificationTokenRepo.GetByHash build: %w", err)
	}

	var t domain.VerificationToken
	var purpose string
	var usedAt sql.NullTime
	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&t.ID, &t.UserID, &purpose, &t.TokenHash, &t.ExpiresAt, &usedAt, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("VerificationTokenRepo.GetByHash scan: %w", err)
	}
	t.Purpose = domain.TokenPurpose(purpose)
	if usedAt.Valid {
		t.UsedAt = &usedAt.Time
	}

	return &t, nil
}

func (r *VerificationTokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Update("verification_tokens").
		Set("used_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("VerificationTokenRepo.MarkUsed build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("VerificationTokenRepo.MarkUsed exec: %w", err)
	}
	return nil
}

func (r *VerificationTokenRepository) InvalidateAllForUser(ctx context.Context, userID uuid.UUID, purpose domain.TokenPurpose) error {
	query, args, err := sq.Update("verification_tokens").
		Set("used_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID, "purpose": string(purpose)}).
		Where("used_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("VerificationTokenRepo.InvalidateAllForUser build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("VerificationTokenRepo.InvalidateAllForUser exec: %w", err)
	}
	return nil
}
