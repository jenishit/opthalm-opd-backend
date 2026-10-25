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

// CreateRole godoc
//
//	@Summary		Create a role
//	@Description	Admin-only. Role names are global (not clinic-scoped) reusable labels.
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateRole	true	"Role name"
//	@Success		200		{object}	response{data=domain.Role}
//	@Failure		409		{object}	errorResponse	"role already exists"
//	@Router			/role/create [post]
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
