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
