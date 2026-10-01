package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
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
//	@Failure		404	{object}	errorResponse
//	@Router			/visit/{id} [get]
func (vh *VisitHandler) GetVisitByVisitID(ctx *gin.Context) {
	visitID := ctx.Param("id")
	visitUUID, err := uuid.Parse(visitID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	visit, err := vh.svc.GetVisitByVisitID(ctx, user.ClinicID, visitUUID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.VisitResponseFromDetails(visit)
	handleSuccess(ctx, rsp)
}

// GetVisitsByPatientID godoc
//
//	@Summary		List a patient's visits
//	@Tags			visits
//	@Produce		json
//	@Security		BearerAuth
//	@Param			patientId	path		string	true	"Patient ID"
//	@Success		200			{object}	response{data=dto.PatientVisitsResponse}
//	@Router			/visit/patient/{patientId} [get]
func (vh *VisitHandler) GetVisitsByPatientID(ctx *gin.Context) {
	patientID := ctx.Param("patientId")
	patientUUID, err := uuid.Parse(patientID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	visits, err := vh.svc.GetVisitsByPatientID(ctx, user.ClinicID, patientUUID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.PatientVisitsResponseFromDetails(visits)
	handleSuccess(ctx, rsp)
}

// UpdateVisitByVisitID godoc
//
//	@Summary		Update a visit
//	@Tags			visits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Visit ID"
//	@Param			request	body		dto.CreateVisitReq	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/visit/{id} [patch]
func (vh *VisitHandler) UpdateVisitByVisitID(ctx *gin.Context) {
	visitID := ctx.Param("id")
	visitUUID, err := uuid.Parse(visitID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

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
		ID:             visitUUID,
		PatientID:      req.PatientID,
		ExamineBy:      req.ExamineBy,
		Status:         req.Status,
		VisitDate:      req.VisitDate,
		CheifComplaint: req.CheifComplaint,
		UpdatedBy:      user.UserId,
	}

	err = vh.svc.UpdateVisitByVisitID(ctx, user.ClinicID, visit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Visit updated successfully"})
}
