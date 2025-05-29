package http

import (
	"bytes"
	"image/png"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

// ─── Inventory Items ─────────────────────────────────────────────

type InventoryHandler struct {
	svc port.InventoryItemService
}

func NewInventoryHandler(svc port.InventoryItemService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// CreateItem godoc
//
//	@Summary		Create an inventory item
//	@Description	Requires ROLE_ADMIN or ROLE_INVENTORY.
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateInventoryItemReq	true	"Item details"
//	@Success		200		{object}	response{data=dto.InventoryItemResponse}
//	@Router			/inventory/items [post]
func (h *InventoryHandler) CreateItem(ctx *gin.Context) {
	var req dto.CreateInventoryItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	unit := req.Unit
	if unit == "" {
		unit = "pcs"
	}

	item := &domain.InventoryItem{
		Category:         domain.InventoryCategory(req.Category),
		SKU:              req.SKU,
		Name:             req.Name,
		Brand:            req.Brand,
		Model:            req.Model,
		Color:            req.Color,
		Size:             req.Size,
		CostPrice:        req.CostPrice,
