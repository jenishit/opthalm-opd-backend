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

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUser

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	user, err := uh.usvc.CreateUser(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, user)
}
