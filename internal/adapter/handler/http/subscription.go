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
