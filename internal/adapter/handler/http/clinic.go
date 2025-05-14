package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type ClinicHandler struct {
	svc port.ClinicService
}

func NewClinicHandler(svc port.ClinicService) *ClinicHandler {
	return &ClinicHandler{
		svc: svc,
	}
}

// InsertClinic godoc
//
//	@Summary		Create clinic settings
//	@Description	Admin-only.
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ClinicRequest	true	"Clinic details"
//	@Success		200		{object}	response{data=dto.ClinicResponse}
