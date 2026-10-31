package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ClinicRepository struct {
	DB *postgres.DB
}

func NewClinicRepository(db *postgres.DB) *ClinicRepository {
	return &ClinicRepository{
		DB: db,
	}
}

func (sr *ClinicRepository) InsertClinic(ctx context.Context, s *domain.ClinicSettings) (*domain.ClinicSettings, error) {
	now := time.Now()

	// require_consent_before_exam/low_stock_notifications are NOT NULL
	// columns, so an omitted field on create needs an explicit default
	// rather than a NULL value.
	requireConsent := false
	if s.RequireConsentBeforeExam != nil {
		requireConsent = *s.RequireConsentBeforeExam
	}
	lowStockNotifications := true
	if s.LowStockNotifications != nil {
		lowStockNotifications = *s.LowStockNotifications
	}

	query, args, err := sq.
		Insert("clinic_settings").
		Columns(
			"clinic_name",
			"tagline",
			"address",
			"phone",
			"email",
			"registration_no",
			"report_footer",
			"workspace_slug",
			"default_currency",
			"timezone",
			"require_consent_before_exam",
			"low_stock_notifications",
			"updated_at",
			"updated_by ",
		).
		Values(
			s.ClinicName,
			s.Tagline,
			s.Address,
			s.Phone,
			s.Email,
			s.RegistrationNo,
			s.ReportFooter,
			s.WorkspaceSlug,
			s.DefaultCurrency,
			s.Timezone,
			requireConsent,
			lowStockNotifications,
			now,
			s.UpdatedBy,
		).
		Suffix(`RETURNING
			id,
			clinic_name
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = sr.DB.QueryRow(ctx, query, args...).Scan(
		&s.ID,
		&s.ClinicName,
	)

	if err != nil {
		return nil, fmt.Errorf("inserting clinic preference: %w", err)
	}

	return s, nil
}

func (sr *ClinicRepository) GetClinicByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicSettings, error) {
	var phone, email, tagline, address, registrationNo, reportFooter sql.NullString
	var workspaceSlug, defaultCurrency, timezone sql.NullString

	query, args, err := sq.
		Select(
			"id",
			"clinic_name",
			"tagline",
			"address",
			"phone",
			"email",
			"registration_no",
			"report_footer",
			"workspace_slug",
			"default_currency",
			"timezone",
			"require_consent_before_exam",
			"low_stock_notifications",
			"updated_by",
		).From("clinic_settings").
		Where(sq.Eq{"id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := sr.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close() //Release resources after reading the rows

	clinic := &domain.ClinicSettings{}
	err = sr.DB.QueryRow(ctx, query, args...).Scan(
		&clinic.ID,
		&clinic.ClinicName,
		&tagline,
		&address,
		&phone,
		&email,
		&registrationNo,
		&reportFooter,
		&workspaceSlug,
		&defaultCurrency,
		&timezone,
		&clinic.RequireConsentBeforeExam,
		&clinic.LowStockNotifications,
		&clinic.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	if workspaceSlug.Valid {
		clinic.WorkspaceSlug = &workspaceSlug.String
	}
	if defaultCurrency.Valid {
		clinic.DefaultCurrency = &defaultCurrency.String
	}
	if timezone.Valid {
		clinic.Timezone = &timezone.String
	}
	if tagline.Valid {
		clinic.Tagline = &tagline.String
	}
	if address.Valid {
		clinic.Address = &address.String
	}
	if email.Valid {
		clinic.Email = &email.String
	}
	if phone.Valid {
		clinic.Phone = &phone.String
	}
	if registrationNo.Valid {
		clinic.RegistrationNo = &registrationNo.String
	}
	if reportFooter.Valid {
		clinic.ReportFooter = &reportFooter.String
	}

	return clinic, nil
}

func (sr *ClinicRepository) UpdateClinic(ctx context.Context, s *domain.ClinicSettings) error {

	builder := sq.Update("clinic_settings").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": s.ID})

	if s.ClinicName != "" {
		builder = builder.Set("clinic_name", s.ClinicName)
	}
	if s.Tagline != nil {
		builder = builder.Set("tagline", s.Tagline)
	}
	if s.Address != nil {
		builder = builder.Set("address", s.Address)
	}
	if s.Phone != nil {
		builder = builder.Set("phone", s.Phone)
	}

	if s.ReportFooter != nil {
		builder = builder.Set("report_footer", s.ReportFooter)
	}
	if s.WorkspaceSlug != nil {
		builder = builder.Set("workspace_slug", s.WorkspaceSlug)
	}
	if s.DefaultCurrency != nil {
		builder = builder.Set("default_currency", s.DefaultCurrency)
	}
	if s.Timezone != nil {
		builder = builder.Set("timezone", s.Timezone)
	}
	if s.RequireConsentBeforeExam != nil {
		builder = builder.Set("require_consent_before_exam", s.RequireConsentBeforeExam)
	}
	if s.LowStockNotifications != nil {
		builder = builder.Set("low_stock_notifications", s.LowStockNotifications)
	}

	builder = builder.Set("updated_by", s.UpdatedBy)

	query, args, err := builder.ToSql()

	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = sr.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

func (sr *ClinicRepository) GetAllClinics(ctx context.Context) ([]*domain.ClinicSettings, error) {
	var phone, email, tagline, address, registrationNo, reportFooter sql.NullString
	var workspaceSlug, defaultCurrency, timezone sql.NullString

	query, args, err := sq.
		Select(
			"id",
			"clinic_name",
			"tagline",
			"address",
			"phone",
			"email",
			"registration_no",
			"report_footer",
			"workspace_slug",
			"default_currency",
			"timezone",
			"require_consent_before_exam",
			"low_stock_notifications",
			"updated_by",
		).From("clinic_settings").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := sr.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying clinics: %w", err)
	}

	defer rows.Close()

	var clinics []*domain.ClinicSettings

	for rows.Next() {
		var clinic domain.ClinicSettings

		err := rows.Scan(
			&clinic.ID,
			&clinic.ClinicName,
			&tagline,
			&address,
			&phone,
			&email,
			&registrationNo,
			&reportFooter,
			&workspaceSlug,
			&defaultCurrency,
			&timezone,
			&clinic.RequireConsentBeforeExam,
			&clinic.LowStockNotifications,
			&clinic.UpdatedBy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if workspaceSlug.Valid {
			clinic.WorkspaceSlug = &workspaceSlug.String
		}
		if defaultCurrency.Valid {
			clinic.DefaultCurrency = &defaultCurrency.String
		}
		if timezone.Valid {
			clinic.Timezone = &timezone.String
		}
		if tagline.Valid {
			clinic.Tagline = &tagline.String
		}
		if address.Valid {
			clinic.Address = &address.String
		}
		if email.Valid {
			clinic.Email = &email.String
		}
		if phone.Valid {
			clinic.Phone = &phone.String
		}
		if registrationNo.Valid {
			clinic.RegistrationNo = &registrationNo.String
		}
		if reportFooter.Valid {
			clinic.ReportFooter = &reportFooter.String
		}

		clinics = append(clinics, &clinic)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating clinic rows: %w", err)
	}

	return clinics, nil
}
