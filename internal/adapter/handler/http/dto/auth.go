package dto

import (
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string     `json:"access_token"`
	SessionID   uuid.UUID  `json:"session_id"`
	UserID      uuid.UUID  `json:"user_id"`
	UserRole    string     `json:"user_role"`
}

func ToLoginResponse(result *domain.LoginResponse) LoginResponse {
	return LoginResponse{
		AccessToken: result.AccessToken,
		SessionID:   result.SessionID,
		UserID:      result.UserID,
		UserRole:    result.UserRole,
	}
}
