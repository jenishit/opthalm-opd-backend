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
			ToSql()
		if err != nil {
			return fmt.Errorf("SessionRepo.Rotate build insert: %w", err)
		}
		if err := tx.QueryRow(ctx, query, args...).Scan(&newSession.ID, &newSession.IssuedAt); err != nil {
			return fmt.Errorf("SessionRepo.Rotate insert: %w", err)
		}

		uQuery, uArgs, err := sq.Update("sessions").
			Set("revoked_at", sq.Expr("NOW()")).
			Set("replaced_by", newSession.ID).
			Where(sq.Eq{"id": oldSessionID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("SessionRepo.Rotate build update: %w", err)
		}
		if _, err := tx.Exec(ctx, uQuery, uArgs...); err != nil {
			return fmt.Errorf("SessionRepo.Rotate update: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newSession, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Update("sessions").
		Set("revoked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Where("revoked_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo.Revoke build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("SessionRepo.Revoke exec: %w", err)
	}
	return nil
}

func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query, args, err := sq.Update("sessions").
		Set("revoked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID}).
		Where("revoked_at IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("SessionRepo.RevokeAllForUser build: %w", err)
	}
	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("SessionRepo.RevokeAllForUser exec: %w", err)
	}
	return nil
}

func scanSession(row pgx.Row) (*domain.Session, error) {
