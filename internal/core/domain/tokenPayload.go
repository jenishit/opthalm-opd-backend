package domain

import (
	"github.com/google/uuid"
)

type TokenPayload struct {
	UserId    uuid.UUID `json:"user_id"`
	RoleName  string    `json:"role_name"`
	ClinicID  uuid.UUID `json:"clinic_id"`
	SessionID uuid.UUID `json:"session_id"`
}
