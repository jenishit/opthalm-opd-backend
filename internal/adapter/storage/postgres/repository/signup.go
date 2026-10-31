package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
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
			 ON CONFLICT (role_name) DO UPDATE SET role_name = EXCLUDED.role_name
			 RETURNING id`,
		).Scan(&roleID); err != nil {
			return fmt.Errorf("SignupRepo.CreateTenant ensure role: %w", err)
		}

		if err := tx.QueryRow(ctx,
			`INSERT INTO users (role_id, clinic_id, email, password) VALUES ($1, $2, $3, $4) RETURNING id`,
			roleID, clinicID, req.AdminEmail, passwordHash,
		).Scan(&userID); err != nil {
			return fmt.Errorf("SignupRepo.CreateTenant insert user: %w", err)
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO profile (user_id, first_name, last_name) VALUES ($1, $2, $3)`,
			userID, req.AdminFirstName, req.AdminLastName,
		); err != nil {
			return fmt.Errorf("SignupRepo.CreateTenant insert profile: %w", err)
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO subscriptions (clinic_id, status, current_period_end) VALUES ($1, 'trialing', now() + interval '14 days')`,
			clinicID,
		); err != nil {
			return fmt.Errorf("SignupRepo.CreateTenant insert subscription: %w", err)
		}

		user.ID = userID
		user.ClinicID = clinicID

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
