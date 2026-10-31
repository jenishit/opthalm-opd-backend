package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

// PatientHistoryHandler serves a patient's clinical history timeline
// (background/PMH): known conditions, current medications, spectacle use.
// It's an append-only log — entries are never edited or deleted, only added.
type PatientHistoryHandler struct {
	svc port.PatientHistoryService
}

func NewPatientHistoryHandler(svc port.PatientHistoryService) *PatientHistoryHandler {
	return &PatientHistoryHandler{svc: svc}
}

// ListPatientHistory godoc
//
//	@Summary		List a patient's clinical history, newest first
//	@Tags			patients
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Patient ID"
//	@Success		200	{object}	response{data=[]dto.PatientHistoryResponse}
//	@Router			/patient/{id}/history [get]
func (h *PatientHistoryHandler) ListPatientHistory(ctx *gin.Context) {
	patientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	history, err := h.svc.ListPatientHistoryByPatientID(ctx, user.ClinicID, patientID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PatientHistoryResList(history))
}

// CreatePatientHistory godoc
//
//	@Summary		Append a clinical history entry
//	@Description	History is a log: this always appends a new entry, never edits an existing one.
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Patient ID"
//	@Param			request	body		dto.CreatePatientHistoryReq	true	"History entry"
//	@Success		200		{object}	response{data=dto.PatientHistoryResponse}
//	@Router			/patient/{id}/history [post]
func (h *PatientHistoryHandler) CreatePatientHistory(ctx *gin.Context) {
	patientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.CreatePatientHistoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	entry := &domain.PatientHistory{
		PatientID:           patientID,
		UsesSpectacles:      req.UsesSpectacles,
		SystemicMedications: req.SystemicMedications,
		TopicalMedications:  req.TopicalMedications,
		Notes:               req.Notes,
		ConditionIDs:        req.ConditionIDs,
		CreatedBy:           user.UserId,
	}

	saved, err := h.svc.CreatePatientHistory(ctx, user.ClinicID, entry)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PatientHistoryRes(saved))
}
