package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PrivacyConfigRepository struct {
	DB *postgres.DB
}

func NewPrivacyConfigRepository(db *postgres.DB) *PrivacyConfigRepository {
	return &PrivacyConfigRepository{DB: db}
}

func (r *PrivacyConfigRepository) UpsertPrivacyConfig(ctx context.Context, c *domain.PrivacyConfig) (*domain.PrivacyConfig, error) {
	query, args, err := sq.
		Insert("clinic_privacy_config").
		Columns(
			"clinic_id", "retention_years", "auto_anonymize_after_retention",
			"allow_patient_data_export", "allow_patient_deletion_requests",
			"require_2fa", "session_timeout_minutes", "login_alerts",
		).
		Values(
			c.ClinicID, c.RetentionYears, c.AutoAnonymizeAfterRetention,
			c.AllowPatientDataExport, c.AllowPatientDeletionRequests,
			c.Require2FA, c.SessionTimeoutMinutes, c.LoginAlerts,
		).
		Suffix(`ON CONFLICT (clinic_id) DO UPDATE SET
			retention_years = EXCLUDED.retention_years,
			auto_anonymize_after_retention = EXCLUDED.auto_anonymize_after_retention,
			allow_patient_data_export = EXCLUDED.allow_patient_data_export,
			allow_patient_deletion_requests = EXCLUDED.allow_patient_deletion_requests,
			require_2fa = EXCLUDED.require_2fa,
			session_timeout_minutes = EXCLUDED.session_timeout_minutes,
			login_alerts = EXCLUDED.login_alerts,
			updated_at = NOW()
		RETURNING clinic_id, retention_years, auto_anonymize_after_retention,
			allow_patient_data_export, allow_patient_deletion_requests,
			require_2fa, session_timeout_minutes, login_alerts`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PrivacyConfig.Upsert build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&c.ClinicID, &c.RetentionYears, &c.AutoAnonymizeAfterRetention,
		&c.AllowPatientDataExport, &c.AllowPatientDeletionRequests,
		&c.Require2FA, &c.SessionTimeoutMinutes, &c.LoginAlerts,
	)
	if err != nil {
		return nil, fmt.Errorf("PrivacyConfig.Upsert scan: %w", err)
	}

	return c, nil
}

func (r *PrivacyConfigRepository) GetPrivacyConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PrivacyConfig, error) {
	query, args, err := sq.
		Select(
			"clinic_id", "retention_years", "auto_anonymize_after_retention",
			"allow_patient_data_export", "allow_patient_deletion_requests",
			"require_2fa", "session_timeout_minutes", "login_alerts",
		).
		From("clinic_privacy_config").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var c domain.PrivacyConfig

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&c.ClinicID, &c.RetentionYears, &c.AutoAnonymizeAfterRetention,
		&c.AllowPatientDataExport, &c.AllowPatientDeletionRequests,
		&c.Require2FA, &c.SessionTimeoutMinutes, &c.LoginAlerts,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.PrivacyConfig{
				ClinicID:               clinicID,
				RetentionYears:         7,
				AllowPatientDataExport: true,
				LoginAlerts:            true,
			}, nil
		}
		return nil, fmt.Errorf("PrivacyConfig.Get scan: %w", err)
	}

	return &c, nil
}

type AuditLogRepository struct {
	DB *postgres.DB
}

func NewAuditLogRepository(db *postgres.DB) *AuditLogRepository {
	return &AuditLogRepository{DB: db}
}

func (r *AuditLogRepository) RecordAuditLogEntry(ctx context.Context, e *domain.AuditLogEntry) (*domain.AuditLogEntry, error) {
	query, args, err := sq.
		Insert("clinic_audit_log").
		Columns("clinic_id", "actor_user_id", "action", "target_type", "target_id").
		Values(e.ClinicID, e.ActorUserID, e.Action, e.TargetType, e.TargetID).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AuditLog.Record build: %w", err)
	}

	if err := r.DB.QueryRow(ctx, query, args...).Scan(&e.ID, &e.CreatedAt); err != nil {
		return nil, fmt.Errorf("AuditLog.Record scan: %w", err)
	}

	return e, nil
}

func (r *AuditLogRepository) ListAuditLog(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.AuditLogEntry, error) {
	query, args, err := sq.
		Select("id", "clinic_id", "actor_user_id", "action", "target_type", "target_id", "created_at").
		From("clinic_audit_log").
		Where(sq.Eq{"clinic_id": clinicID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("AuditLog.List query: %w", err)
	}
	defer rows.Close()

	entries := make([]*domain.AuditLogEntry, 0)

	for rows.Next() {
		var e domain.AuditLogEntry

		if err := rows.Scan(&e.ID, &e.ClinicID, &e.ActorUserID, &e.Action, &e.TargetType, &e.TargetID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("AuditLog.List scan: %w", err)
		}

		entries = append(entries, &e)
	}

	return entries, nil
}
