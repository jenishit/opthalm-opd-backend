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
		SellingPrice:     req.SellingPrice,
		QuantityOnHand:   req.QuantityOnHand,
		ReorderThreshold: req.ReorderThreshold,
		Unit:             unit,
		CreatedBy:        user.UserId,
		UpdatedBy:        user.UserId,
	}

	created, err := h.svc.Create(ctx, user.ClinicID, item)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InventoryItemRes(created))
}

// ListItems godoc
//
//	@Summary		List inventory items
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.InventoryItemResponse}
//	@Router			/inventory/items [get]
func (h *InventoryHandler) ListItems(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	items, err := h.svc.List(ctx, user.ClinicID, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InventoryItemResList(items))
}

// SearchItems godoc
//
//	@Summary		Search inventory items
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.InventoryItemResponse}
//	@Router			/inventory/items/search [get]
func (h *InventoryHandler) SearchItems(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	items, err := h.svc.Search(ctx, user.ClinicID, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InventoryItemResList(items))
}

// GetItemByID godoc
//
//	@Summary		Get an inventory item
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Item ID"
//	@Success		200	{object}	response{data=dto.InventoryItemResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/inventory/items/{id} [get]
func (h *InventoryHandler) GetItemByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	item, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InventoryItemRes(item))
}

// GetItemByBarcode godoc
//
//	@Summary		Get an inventory item by SKU/barcode
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			sku	path		string	true	"SKU"
//	@Success		200	{object}	response{data=dto.InventoryItemResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/inventory/items/barcode/{sku} [get]
func (h *InventoryHandler) GetItemByBarcode(ctx *gin.Context) {
	sku := ctx.Param("sku")

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	item, err := h.svc.GetBySKU(ctx, user.ClinicID, sku)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InventoryItemRes(item))
}

// UpdateItem godoc
//
//	@Summary		Update an inventory item
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Item ID"
//	@Param			request	body		dto.UpdateInventoryItemReq	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/inventory/items/{id} [patch]
func (h *InventoryHandler) UpdateItem(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateInventoryItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
