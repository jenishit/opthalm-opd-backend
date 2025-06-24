package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type PatientHandler struct {
	svc port.PatientService
}

func NewPatientHandler(svc port.PatientService) *PatientHandler {
	return &PatientHandler{
		svc: svc,
	}
}

// CreatePatient godoc
//
//	@Summary		Create a patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreatePatientReq	true	"Patient details"
//	@Success		200		{object}	response{data=dto.PatientResponse}
//	@Router			/patient [post]
func (ph *PatientHandler) CreatePatient(ctx *gin.Context) {
	var req dto.CreatePatientReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	patient := &domain.Patient{
		FullName:   req.FullName,
		DOB:        req.DOB,
		Gender:     req.Gender,
		Phone:      req.Phone,
		Occupation: req.Occupation,
		Email:      req.Email,
		Address:    req.Address,
		CreatedBy:  user.UserId,
	}

	p, err := ph.svc.CreatePatient(ctx, user.ClinicID, patient)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PatientRes(p))
}

// GetPatientByID godoc
//
//	@Summary		Get a patient
//	@Tags			patients
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Patient ID"
//	@Success		200	{object}	response{data=dto.PatientResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/patient/{id} [get]
func (ph *PatientHandler) GetPatientByID(ctx *gin.Context) {
	ptID := ctx.Param("id")
	ptUUID, err := uuid.Parse(ptID)

	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

