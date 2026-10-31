package domain

import "github.com/google/uuid"

type NotificationPrefs struct {
	ClinicID                uuid.UUID
	EmailNotifications      bool
	SMSAppointmentReminders bool
	LowStockAlerts          bool
	NewOrderNotifications   bool
	WeeklyDigest            bool
}
