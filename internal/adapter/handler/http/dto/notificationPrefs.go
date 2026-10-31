package dto

import (
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type UpsertNotificationPrefsReq struct {
	EmailNotifications      bool `json:"email_notifications"`
	SMSAppointmentReminders bool `json:"sms_appointment_reminders"`
	LowStockAlerts          bool `json:"low_stock_alerts"`
	NewOrderNotifications   bool `json:"new_order_notifications"`
	WeeklyDigest            bool `json:"weekly_digest"`
}

type NotificationPrefsResponse struct {
	ClinicID                uuid.UUID `json:"clinic_id"`
	EmailNotifications      bool      `json:"email_notifications"`
	SMSAppointmentReminders bool      `json:"sms_appointment_reminders"`
	LowStockAlerts          bool      `json:"low_stock_alerts"`
	NewOrderNotifications   bool      `json:"new_order_notifications"`
	WeeklyDigest            bool      `json:"weekly_digest"`
}

func NotificationPrefsRes(p *domain.NotificationPrefs) *NotificationPrefsResponse {
	if p == nil {
		return nil
	}

	return &NotificationPrefsResponse{
		ClinicID:                p.ClinicID,
		EmailNotifications:      p.EmailNotifications,
		SMSAppointmentReminders: p.SMSAppointmentReminders,
		LowStockAlerts:          p.LowStockAlerts,
		NewOrderNotifications:   p.NewOrderNotifications,
		WeeklyDigest:            p.WeeklyDigest,
	}
}
