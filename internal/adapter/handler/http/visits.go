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
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	visit := &domain.Visit{
		PatientID:      req.PatientID,
		ExamineBy:      req.ExamineBy,
		Status:         req.Status,
		VisitDate:      req.VisitDate,
		CheifComplaint: req.CheifComplaint,
		CreatedBy:      user.UserId,
		UpdatedBy:      user.UserId,
	}

	v, err := vh.svc.CreateVisit(ctx, user.ClinicID, visit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VisitCreateRes(v))
}

// GetVisitByVisitID godoc
//
//	@Summary		Get a visit
//	@Tags			visits
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Visit ID"
//	@Success		200	{object}	response{data=dto.VisitResponse}
