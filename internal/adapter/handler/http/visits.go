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

func (vh *VisitHandler) CreateVisit(ctx *gin.Context) {
	var req dto.CreateVisitReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		validationError(ctx, domain.ErrEmptyAuthorizationHeader)
		return
	}

	userPayload, ok := payload.(*domain.TokenPayload)
	if !ok {
		validationError(ctx, domain.ErrInvalidAuthorizationHeader)
		return
	}

	visit := &domain.Visit{
		PatientID:      req.PatientID,
		ExamineBy:      req.ExamineBy,
		Status:         req.Status,
		VisitDate:      req.VisitDate,
		CheifComplaint: req.CheifComplaint,
		CreatedBy:      userPayload.UserId,
		UpdatedBy:      userPayload.UserId,
	}

	v, err := vh.svc.CreateVisit(ctx, visit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, v)
}

func (vh *VisitHandler) GetVisitByVisitID(ctx *gin.Context) {
	visitID := ctx.Param("id")
	visitUUID, err := uuid.Parse(visitID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	visit, err := vh.svc.GetVisitByVisitID(ctx, visitUUID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.VisitResponseFromDetails(visit)
	handleSuccess(ctx, rsp)
}

func (vh *VisitHandler) GetVisitsByPatientID(ctx *gin.Context) {
	patientID := ctx.Param("patientId")
	patientUUID, err := uuid.Parse(patientID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	visits, err := vh.svc.GetVisitsByPatientID(ctx, patientUUID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.PatientVisitsResponseFromDetails(visits)
	handleSuccess(ctx, rsp)
}

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

	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		validationError(ctx, domain.ErrEmptyAuthorizationHeader)
		return
	}

	userPayload, ok := payload.(*domain.TokenPayload)
	if !ok {
		validationError(ctx, domain.ErrInvalidAuthorizationHeader)
		return
	}

	visit := &domain.Visit{
		ID:             visitUUID,
		PatientID:      req.PatientID,
		ExamineBy:      req.ExamineBy,
		Status:         req.Status,
		VisitDate:      req.VisitDate,
		CheifComplaint: req.CheifComplaint,
		UpdatedBy:      userPayload.UserId,
	}

	err = vh.svc.UpdateVisitByVisitID(ctx, visit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Visit updated successfully"})
}
