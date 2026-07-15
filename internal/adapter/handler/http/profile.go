package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type ProfileHandler struct {
	psvc port.ProfileService
}

func NewProfileHandler(psvc port.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		psvc: psvc,
	}
}

func (ph *ProfileHandler) GetProfileByID(ctx *gin.Context) {
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

	res, err := ph.psvc.GetProfileByID(ctx, userPayload.UserId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp := dto.NewProfileResponse(res)

	handleSuccess(ctx, rsp)
}

func (ph *ProfileHandler) GetProfiles(ctx *gin.Context) {
	res, err := ph.psvc.GetProfiles(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp := dto.NewProfileResponses(res)

	handleSuccess(ctx, rsp)
}

func (ph *ProfileHandler) UpdateProfileByUserID(ctx *gin.Context) {
	user_id := ctx.Param("id")
	id, err := uuid.Parse(user_id)
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	profile := &domain.GetProfileDetails{
		UserID:    id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	err = ph.psvc.UpdateProfileByUserID(ctx, profile)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, "profile updated successfully")
}
