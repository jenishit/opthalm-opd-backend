package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type RoleHandler struct {
	rsvc port.RoleService
}

func NewRoleHandler(rsvc port.RoleService) *RoleHandler {
	return &RoleHandler{
		rsvc: rsvc,
	}
}

func (rh *RoleHandler) CreateRole(ctx *gin.Context) {
	var req dto.CreateRole

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	role := &domain.Role{
		RoleName: req.RoleName,
	}

	role, err := rh.rsvc.CreateRole(ctx, role)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, role)
}
