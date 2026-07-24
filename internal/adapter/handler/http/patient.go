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

func (ph *PatientHandler) CreatePatient(ctx *gin.Context) {
	var req dto.CreatePatientReq

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

	patient := &domain.Patient{
		FullName:   req.FullName,
		DOB:        req.DOB,
		Gender:     req.Gender,
		Phone:      req.Phone,
		Occupation: req.Occupation,
		Email:      req.Email,
		Address:    req.Address,
		CreatedBy:  userPayload.UserId,
	}

	p, err := ph.svc.CreatePatient(ctx, patient)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, p)
}

func (ph *PatientHandler) GetPatientByID(ctx *gin.Context) {
	ptID := ctx.Param("id")
	ptUUID, err := uuid.Parse(ptID)

	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	patient, err := ph.svc.GetPatientByID(ctx, ptUUID)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, patient)
}

func (ph *PatientHandler) GetPatients(ctx *gin.Context) {
	res, err := ph.svc.GetPatients(ctx)

	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.PatientResponses(res)

	handleSuccess(ctx, rsp)
}

func (ph *PatientHandler) SearchPatients(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	patients, err := ph.svc.SearchPatients(ctx, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.PatientResponses(patients)
	handleSuccess(ctx, rsp)
}

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

	patient := &domain.Patient{
		ID:         ptUUID,
		FullName:   req.FullName,
		DOB:        req.DOB,
		Gender:     req.Gender,
		Phone:      req.Phone,
		Occupation: req.Occupation,
		Email:      req.Email,
		Address:    req.Address,
		UpdatedBy:  userPayload.UserId,
	}

	err = ph.svc.UpdatePatientByID(ctx, patient)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Patient updated successfully"})
}

func (ph *PatientHandler) DeletePatientByID(ctx *gin.Context) {
	ptID := ctx.Param("id")
	ptUUID, err := uuid.Parse(ptID)

	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	err = ph.svc.DeletePatientByID(ctx, ptUUID)

	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Patient deleted successfully"})
}
