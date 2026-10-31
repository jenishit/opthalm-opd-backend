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

type NotificationPrefsRepository struct {
	DB *postgres.DB
}

func NewNotificationPrefsRepository(db *postgres.DB) *NotificationPrefsRepository {
	return &NotificationPrefsRepository{DB: db}
}

func (r *NotificationPrefsRepository) UpsertNotificationPrefs(ctx context.Context, p *domain.NotificationPrefs) (*domain.NotificationPrefs, error) {
	query, args, err := sq.
		Insert("clinic_notification_prefs").
		Columns(
			"clinic_id",
			"email_notifications",
			"sms_appointment_reminders",
			"low_stock_alerts",
			"new_order_notifications",
			"weekly_digest",
		).
		Values(
			p.ClinicID,
			p.EmailNotifications,
			p.SMSAppointmentReminders,
			p.LowStockAlerts,
			p.NewOrderNotifications,
			p.WeeklyDigest,
		).
		Suffix(`ON CONFLICT (clinic_id) DO UPDATE SET
			email_notifications = EXCLUDED.email_notifications,
			sms_appointment_reminders = EXCLUDED.sms_appointment_reminders,
			low_stock_alerts = EXCLUDED.low_stock_alerts,
			new_order_notifications = EXCLUDED.new_order_notifications,
			weekly_digest = EXCLUDED.weekly_digest,
			updated_at = NOW()
		RETURNING clinic_id, email_notifications, sms_appointment_reminders, low_stock_alerts, new_order_notifications, weekly_digest`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("NotificationPrefs.Upsert build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&p.ClinicID,
		&p.EmailNotifications,
		&p.SMSAppointmentReminders,
		&p.LowStockAlerts,
		&p.NewOrderNotifications,
		&p.WeeklyDigest,
	)
	if err != nil {
		return nil, fmt.Errorf("NotificationPrefs.Upsert scan: %w", err)
	}

	return p, nil
}

// GetNotificationPrefs returns the clinic's defaults (matching the column
// defaults) when no row has ever been saved, rather than ErrDataNotFound —
// notification prefs always have a sensible value to show.
func (r *NotificationPrefsRepository) GetNotificationPrefs(ctx context.Context, clinicID uuid.UUID) (*domain.NotificationPrefs, error) {
	query, args, err := sq.
		Select(
			"clinic_id",
			"email_notifications",
			"sms_appointment_reminders",
			"low_stock_alerts",
			"new_order_notifications",
			"weekly_digest",
		).
		From("clinic_notification_prefs").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p domain.NotificationPrefs

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&p.ClinicID,
		&p.EmailNotifications,
		&p.SMSAppointmentReminders,
		&p.LowStockAlerts,
		&p.NewOrderNotifications,
		&p.WeeklyDigest,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.NotificationPrefs{
				ClinicID:                clinicID,
				EmailNotifications:      true,
				SMSAppointmentReminders: false,
				LowStockAlerts:          true,
				NewOrderNotifications:   true,
				WeeklyDigest:            false,
			}, nil
		}
		return nil, fmt.Errorf("NotificationPrefs.Get scan: %w", err)
	}

	return &p, nil
}
