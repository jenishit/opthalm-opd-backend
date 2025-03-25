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

type ProfileRepository struct {
	DB *postgres.DB
}

func NewProfileRepository(db *postgres.DB) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (p *ProfileRepository) CreateProfile(ctx context.Context, pr *domain.Profile) (*domain.Profile, error) {
	query, args, err := sq.
		Insert("profile").
		Columns(
			"user_id",
			"first_name",
			"last_name",
			"phone",
		).
		Values(
			pr.UserID,
			pr.FirstName,
			pr.LastName,
			pr.Phone,
		).
		Suffix(`
			RETURNING
				id,
				created_at,
				updated_at
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ProfileRepo.CreateProfile build: %w", err)
	}

	err = p.DB.QueryRow(ctx, query, args...).Scan(
		&pr.ID,
		&pr.CreatedAt,
		&pr.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ProfileRepo.CreateProfile scan: %w", err)
	}

	return pr, nil
}

func (p *ProfileRepository) GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.GetProfileDetails, error) {
	var phone, email sql.NullString

	query, args, err := sq.
		Select(
			"P.ID",
			"P.FIRST_NAME",
			"P.LAST_NAME",
			"R.ROLE_NAME",
			"P.USER_ID",
			"U.EMAIL",
			"P.PHONE",
		).
		From("PROFILE P").
		LeftJoin("USERS U ON U.ID = P.USER_ID").
		LeftJoin("ROLE R ON R.ID = U.ROLE_ID").
		Where(sq.Eq{"P.USER_ID": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var profile domain.GetProfileDetails

	err = p.DB.QueryRow(ctx, query, args...).Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.LastName,
		&profile.RoleName,
		&profile.UserID,
		&email,
		&phone,
	)

