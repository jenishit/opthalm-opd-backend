package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
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
//	@Param			request	body		dto.CreateUser	true	"New user details"
//	@Success		200		{object}	response{data=domain.User}
//	@Failure		409		{object}	errorResponse	"email already exists"
//	@Router			/user/create [post]
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUser

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	caller, ok := currentUser(ctx)
	if !ok {
		return
	}

	user, err := uh.usvc.CreateUser(ctx, &req, caller.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, user)
}
