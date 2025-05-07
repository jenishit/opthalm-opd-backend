package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/services"
)

// CalculatorHandler exposes the stateless optical math helpers. It depends
// directly on the concrete CalculatorService rather than a port interface
// since there's no I/O to abstract over (see services.CalculatorService).
type CalculatorHandler struct {
	svc *services.CalculatorService
}

func NewCalculatorHandler(svc *services.CalculatorService) *CalculatorHandler {
	return &CalculatorHandler{svc: svc}
}

// Transposition godoc
//
//	@Summary		Transpose a cylinder between plus/minus cylinder notation
//	@Tags			calculators
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.TranspositionReq	true	"Sphere/cylinder/axis"
//	@Success		200		{object}	response{data=dto.TranspositionRes}
//	@Router			/calculators/transposition [post]
func (h *CalculatorHandler) Transposition(ctx *gin.Context) {
	var req dto.TranspositionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	sphere, cylinder, axis := h.svc.TransposeCylinder(req.Sphere, req.Cylinder, req.Axis)
	handleSuccess(ctx, dto.TranspositionRes{Sphere: sphere, Cylinder: cylinder, Axis: axis})
}

// SphericalEquivalent godoc
//
//	@Summary		Compute spherical equivalent
//	@Tags			calculators
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.SphericalEquivalentReq	true	"Sphere/cylinder"
//	@Success		200		{object}	response{data=dto.ScalarResultRes}
//	@Router			/calculators/spherical-equivalent [post]
func (h *CalculatorHandler) SphericalEquivalent(ctx *gin.Context) {
	var req dto.SphericalEquivalentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.ScalarResultRes{Result: h.svc.SphericalEquivalent(req.Sphere, req.Cylinder)})
}

