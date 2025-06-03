package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type LabJobHandler struct {
	svc port.LabJobService
}

func NewLabJobHandler(svc port.LabJobService) *LabJobHandler {
	return &LabJobHandler{svc: svc}
}

// CreateLabJob godoc
//
//	@Summary		Create a lab job
//	@Description	Requires ROLE_ADMIN, ROLE_LAB, or ROLE_BILLING.
//	@Tags			lab-jobs
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateLabJobReq	true	"Lab job details"
//	@Success		200		{object}	response{data=dto.LabJobResponse}
//	@Router			/lab-jobs [post]
func (h *LabJobHandler) CreateLabJob(ctx *gin.Context) {
	var req dto.CreateLabJobReq
