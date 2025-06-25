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

	patient, err := ph.svc.GetPatientByID(ctx, user.ClinicID, ptUUID)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PatientRes(patient))
}

// GetPatients godoc
//
//	@Summary		List patients
//	@Tags			patients
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=[]dto.PatientResponse}
//	@Router			/patient [get]
func (ph *PatientHandler) GetPatients(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	res, err := ph.svc.GetPatients(ctx, user.ClinicID)

	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.PatientResponses(res)

	handleSuccess(ctx, rsp)
}

// SearchPatients godoc
//
//	@Summary		Search patients
//	@Tags			patients
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.PatientResponse}
//	@Router			/patient/search [get]
func (ph *PatientHandler) SearchPatients(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	patients, err := ph.svc.SearchPatients(ctx, user.ClinicID, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.PatientResponses(patients)
	handleSuccess(ctx, rsp)
}

// UpdatePatientByID godoc
//
//	@Summary		Update a patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Patient ID"
//	@Param			request	body		dto.CreatePatientReq	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/patient/{id} [patch]
func (ph *PatientHandler) UpdatePatientByID(ctx *gin.Context) {
	ptID := ctx.Param("id")
	ptUUID, err := uuid.Parse(ptID)

	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

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
		ID:         ptUUID,
		FullName:   req.FullName,
		DOB:        req.DOB,
		Gender:     req.Gender,
		Phone:      req.Phone,
		Occupation: req.Occupation,
		Email:      req.Email,
		Address:    req.Address,
		UpdatedBy:  user.UserId,
	}

	err = ph.svc.UpdatePatientByID(ctx, user.ClinicID, patient)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Patient updated successfully"})
}

// DeletePatientByID godoc
//
//	@Summary		Soft-delete a patient
//	@Tags			patients
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Patient ID"
//	@Success		200	{object}	response
//	@Router			/patient/{id}/delete [patch]
func (ph *PatientHandler) DeletePatientByID(ctx *gin.Context) {
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

	err = ph.svc.DeletePatientByID(ctx, user.ClinicID, ptUUID)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Patient deleted successfully"})
}
