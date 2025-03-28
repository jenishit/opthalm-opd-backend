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

type SessionRepository struct {
	DB *postgres.DB
}

func NewSessionRepository(db *postgres.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *domain.Session) (*domain.Session, error) {
	query, args, err := sq.Insert("sessions").
		Columns("user_id", "refresh_token_hash", "user_agent", "ip_address", "expires_at").
		Values(s.UserID, s.RefreshTokenHash, nullStringPtr(s.UserAgent), nullStringPtr(s.IPAddress), s.ExpiresAt).
		Suffix("RETURNING id, issued_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SessionRepo.Create build: %w", err)
	}
	if err := r.DB.QueryRow(ctx, query, args...).Scan(&s.ID, &s.IssuedAt); err != nil {
		return nil, fmt.Errorf("SessionRepo.Create exec: %w", err)
	}
	return s, nil
}

func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	query, args, err := sq.Select(
		"id", "user_id", "refresh_token_hash", "user_agent", "ip_address",
		"issued_at", "expires_at", "revoked_at", "replaced_by",
	).
		From("sessions").
		Where(sq.Eq{"refresh_token_hash": hash}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SessionRepo.GetByRefreshTokenHash build: %w", err)
	}
	return scanSession(r.DB.QueryRow(ctx, query, args...))
}

func (r *SessionRepository) Rotate(ctx context.Context, oldSessionID uuid.UUID, newSession *domain.Session) (*domain.Session, error) {
	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		query, args, err := sq.Insert("sessions").
			Columns("user_id", "refresh_token_hash", "user_agent", "ip_address", "expires_at").
			Values(newSession.UserID, newSession.RefreshTokenHash, nullStringPtr(newSession.UserAgent), nullStringPtr(newSession.IPAddress), newSession.ExpiresAt).
			Suffix("RETURNING id, issued_at").
			PlaceholderFormat(sq.Dollar).
