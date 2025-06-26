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
