package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type VisitHandler struct {
	svc port.VisitsService
}

func NewVisitHandler(svc port.VisitsService) *VisitHandler {
	return &VisitHandler{
		svc: svc,
	}
}

// CreateVisit godoc
//
//	@Summary		Create a visit
//	@Tags			visits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateVisitReq	true	"Visit details"
//	@Success		200		{object}	response{data=dto.VisitCreateResponse}
//	@Router			/visit [post]
func (vh *VisitHandler) CreateVisit(ctx *gin.Context) {
	var req dto.CreateVisitReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
