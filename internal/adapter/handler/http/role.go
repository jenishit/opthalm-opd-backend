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
