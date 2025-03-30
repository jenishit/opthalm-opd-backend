package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type SignupRepository struct {
	DB *postgres.DB
}

func NewSignupRepository(db *postgres.DB) *SignupRepository {
	return &SignupRepository{DB: db}
}

func (r *SignupRepository) CreateTenant(ctx context.Context, req *domain.SignupRequest, passwordHash string) (*domain.BasicDetails, error) {
	user := &domain.BasicDetails{
		UserRole: "ROLE_ADMIN",
		Email:    &req.AdminEmail,
	}

	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		var clinicID, roleID, userID uuid.UUID

		if err := tx.QueryRow(ctx,
			`INSERT INTO clinic_settings (clinic_name, registration_no) VALUES ($1, $2) RETURNING id`,
			req.ClinicName, req.RegistrationNo,
		).Scan(&clinicID); err != nil {
			return fmt.Errorf("SignupRepo.CreateTenant insert clinic: %w", err)
		}

		if err := tx.QueryRow(ctx,
			`INSERT INTO role (role_name) VALUES ('ROLE_ADMIN')
