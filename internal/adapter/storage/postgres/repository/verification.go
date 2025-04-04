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
