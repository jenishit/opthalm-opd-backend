package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

// currentUser extracts the authenticated TokenPayload set by authMiddleware,
// writing the appropriate error response itself on failure.
func currentUser(ctx *gin.Context) (*domain.TokenPayload, bool) {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		validationError(ctx, domain.ErrEmptyAuthorizationHeader)
		return nil, false
	}

	userPayload, ok := payload.(*domain.TokenPayload)
	if !ok {
		validationError(ctx, domain.ErrInvalidAuthorizationHeader)
		return nil, false
	}

	return userPayload, true
}
