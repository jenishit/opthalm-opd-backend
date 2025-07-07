package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type UserHandler struct {
	usvc port.UserService
}

func NewUsersHandler(usvc port.UserService) *UserHandler {
	return &UserHandler{
		usvc: usvc,
	}
}

// CreateUser godoc
//
//	@Summary		Create a staff user
//	@Description	Admin-only. Creates a new user (and profile) in the caller's own clinic.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
