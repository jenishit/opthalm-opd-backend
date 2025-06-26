package http

import (
	"strconv"

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

// GetProfileByID godoc
//
//	@Summary		Get my profile
//	@Description	Returns the authenticated caller's own profile.
//	@Tags			profile
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=dto.ProfileResponse}
//	@Router			/profile/getme [get]
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

// GetProfiles godoc
//
//	@Summary		List profiles
//	@Tags			profile
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=[]dto.ProfileResponse}
//	@Router			/profile/profile-details [get]
func (ph *ProfileHandler) GetProfiles(ctx *gin.Context) {
	res, err := ph.psvc.GetProfiles(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp := dto.NewProfileResponses(res)

	handleSuccess(ctx, rsp)
}

// SearchProfiles godoc
//
//	@Summary		Search profiles
//	@Tags			profile
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.ProfileResponse}
//	@Router			/profile/search [get]
func (ph *ProfileHandler) SearchProfiles(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	profiles, err := ph.psvc.SearchProfiles(ctx, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := dto.NewProfileResponses(profiles)
	handleSuccess(ctx, rsp)
}

// UpdateProfileByUserID godoc
//
//	@Summary		Update a profile
//	@Tags			profile
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"User ID"
//	@Param			request	body		dto.UpdateProfileRequest	true	"Fields to update"
//	@Success		200		{object}	response
//	@Failure		400		{object}	errorResponse	"invalid UUID"
//	@Router			/profile/update-profile/{id} [patch]
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
