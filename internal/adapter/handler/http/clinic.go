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
//	@Router			/admin/clinic [post]
func (ch *ClinicHandler) InsertClinic(ctx *gin.Context) {
	var req dto.ClinicRequest

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

	clinic := &domain.ClinicSettings{
		ClinicName:     req.ClinicName,
		Tagline:        &req.Tagline,
		Address:        &req.Address,
		Phone:          &req.Phone,
		Email:          &req.Email,
		RegistrationNo: &req.RegistrationNo,
		ReportFooter:   &req.ReportFooter,
		UpdatedBy:      userPayload.UserId,
	}

	s, err := ch.svc.InsertClinic(ctx, clinic)

	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, dto.ClinicsResponse(s))
}

// GetClinicByID godoc
//
//	@Summary		Get clinic settings
//	@Description	Admin-only.
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Clinic ID"
//	@Success		200	{object}	response{data=dto.ClinicResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/admin/clinic/{id} [get]
func (ch *ClinicHandler) GetClinicByID(ctx *gin.Context) {
	clinicID := ctx.Param("id")
	clinicUUID, err := uuid.Parse(clinicID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	clinic, err := ch.svc.GetClinicByClinicID(ctx, clinicUUID)

	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, dto.ClinicsResponse(clinic))
}

// UpdateClinic godoc
//
//	@Summary		Update clinic settings
//	@Description	Admin-only.
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Clinic ID"
//	@Param			request	body		dto.ClinicRequest	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/admin/clinic/{id} [patch]
func (ch *ClinicHandler) UpdateClinic(ctx *gin.Context) {
	clinicID := ctx.Param("id")
	clinicUUID, err := uuid.Parse(clinicID)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.ClinicRequest
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

	clinic := &domain.ClinicSettings{
		ID:             clinicUUID,
		ClinicName:     req.ClinicName,
		Tagline:        &req.Tagline,
		Address:        &req.Address,
		Phone:          &req.Phone,
		Email:          &req.Email,
		RegistrationNo: &req.RegistrationNo,
		ReportFooter:   &req.ReportFooter,
		UpdatedBy:      userPayload.UserId,
	}

