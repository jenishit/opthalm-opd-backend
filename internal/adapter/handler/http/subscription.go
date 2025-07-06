package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

// SubscriptionHandler is platform-operator-only (see superadminMiddleware) —
// a clinic's own admin can never reach these routes, so it can never
// reactivate or extend its own subscription.
type SubscriptionHandler struct {
	svc port.SubscriptionService
}

func NewSubscriptionHandler(svc port.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

// GetByClinicID godoc
//
//	@Summary		Get a clinic's subscription
//	@Description	Platform-operator-only (ROLE_SUPERADMIN).
//	@Tags			platform
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clinicId	path		string	true	"Clinic ID"
//	@Success		200			{object}	response{data=dto.SubscriptionResponse}
//	@Failure		404			{object}	errorResponse
//	@Router			/platform/subscriptions/{clinicId} [get]
func (h *SubscriptionHandler) GetByClinicID(ctx *gin.Context) {
	clinicID, err := uuid.Parse(ctx.Param("clinicId"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	sub, err := h.svc.GetByClinicID(ctx, clinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.SubscriptionRes(sub))
}

// Upsert godoc
//
//	@Summary		Create or update a clinic's subscription
//	@Description	Platform-operator-only (ROLE_SUPERADMIN). A clinic's own admin can never reach this route, so it can never reactivate or extend its own subscription.
//	@Tags			platform
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clinicId	path		string							true	"Clinic ID"
//	@Param			request		body		dto.UpsertSubscriptionReq	true	"Plan, status, period end"
//	@Success		200			{object}	response{data=dto.SubscriptionResponse}
//	@Router			/platform/subscriptions/{clinicId} [put]
func (h *SubscriptionHandler) Upsert(ctx *gin.Context) {
	clinicID, err := uuid.Parse(ctx.Param("clinicId"))
