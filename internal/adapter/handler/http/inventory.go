package http

import (
	"bytes"
	"image/png"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
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
	}

	item := &domain.InventoryItem{
		ID:        id,
		Brand:     req.Brand,
		Model:     req.Model,
		Color:     req.Color,
		Size:      req.Size,
		UpdatedBy: user.UserId,
	}
	if req.Category != nil {
		item.Category = domain.InventoryCategory(*req.Category)
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.CostPrice != nil {
		item.CostPrice = *req.CostPrice
	}
	if req.SellingPrice != nil {
		item.SellingPrice = *req.SellingPrice
	}
	if req.ReorderThreshold != nil {
		item.ReorderThreshold = *req.ReorderThreshold
	}

	if err := h.svc.Update(ctx, user.ClinicID, item); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Inventory item updated successfully"})
}

// DeleteItem godoc
//
//	@Summary		Soft-delete an inventory item
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Item ID"
//	@Success		200	{object}	response
//	@Router			/inventory/items/{id}/delete [patch]
func (h *InventoryHandler) DeleteItem(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	if err := h.svc.Delete(ctx, user.ClinicID, id); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Inventory item deleted successfully"})
}

// LowStock godoc
//
//	@Summary		List items at or below their reorder threshold
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=[]dto.InventoryItemResponse}
//	@Router			/inventory/low-stock [get]
func (h *InventoryHandler) LowStock(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	items, err := h.svc.LowStock(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InventoryItemResList(items))
}

// AddStock godoc
//
//	@Summary		Add stock to an item
//	@Description	Records a stock movement and increments quantity on hand.
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Item ID"
//	@Param			request	body		dto.AddStockReq	true	"Quantity + notes"
//	@Success		200		{object}	response{data=dto.StockMovementResponse}
//	@Router			/inventory/items/{id}/stock [post]
func (h *InventoryHandler) AddStock(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.AddStockReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	movement, err := h.svc.AddStock(ctx, user.ClinicID, id, req.Quantity, req.Notes, user.UserId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.StockMovementRes(movement))
}

// ListMovements godoc
//
//	@Summary		List an item's stock movements
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"Item ID"
//	@Param			limit	query		int		false	"Max results"	default(20)
//	@Param			offset	query		int		false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.StockMovementResponse}
//	@Router			/inventory/items/{id}/movements [get]
func (h *InventoryHandler) ListMovements(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	movements, err := h.svc.ListMovements(ctx, user.ClinicID, id, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.StockMovementResList(movements))
}

// GetBarcodeImage godoc
//
//	@Summary		Get an item's barcode as a PNG image
//	@Tags			inventory
//	@Produce		image/png
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Item ID"
//	@Success		200	{file}	binary
//	@Router			/inventory/items/{id}/barcode-image [get]
func (h *InventoryHandler) GetBarcodeImage(ctx *gin.Context) {
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

	bc, err := code128.Encode(item.SKU)
	if err != nil {
		handleError(ctx, domain.ErrInternal)
		return
	}

	scaled, err := barcode.Scale(bc, 300, 100)
	if err != nil {
		handleError(ctx, domain.ErrInternal)
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		handleError(ctx, domain.ErrInternal)
		return
	}

	ctx.Data(200, "image/png", buf.Bytes())
}

// ─── Vendors ─────────────────────────────────────────────────────

type VendorHandler struct {
	svc port.VendorService
}

func NewVendorHandler(svc port.VendorService) *VendorHandler {
	return &VendorHandler{svc: svc}
}

// CreateVendor godoc
//
//	@Summary		Create a vendor
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateVendorReq	true	"Vendor details"
//	@Success		200		{object}	response{data=dto.VendorResponse}
//	@Router			/inventory/vendors [post]
func (h *VendorHandler) CreateVendor(ctx *gin.Context) {
	var req dto.CreateVendorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	vendor := &domain.Vendor{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
	}

	created, err := h.svc.Create(ctx, user.ClinicID, vendor)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VendorRes(created))
}

// ListVendors godoc
//
//	@Summary		List vendors
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.VendorResponse}
//	@Router			/inventory/vendors [get]
func (h *VendorHandler) ListVendors(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	vendors, err := h.svc.List(ctx, user.ClinicID, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VendorResList(vendors))
}

// SearchVendors godoc
//
//	@Summary		Search vendors
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.VendorResponse}
//	@Router			/inventory/vendors/search [get]
func (h *VendorHandler) SearchVendors(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	vendors, err := h.svc.Search(ctx, user.ClinicID, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VendorResList(vendors))
}

// GetVendorByID godoc
//
//	@Summary		Get a vendor
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Vendor ID"
//	@Success		200	{object}	response{data=dto.VendorResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/inventory/vendors/{id} [get]
func (h *VendorHandler) GetVendorByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	vendor, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VendorRes(vendor))
}

// UpdateVendor godoc
//
//	@Summary		Update a vendor
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Vendor ID"
//	@Param			request	body		dto.UpdateVendorReq	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/inventory/vendors/{id} [patch]
func (h *VendorHandler) UpdateVendor(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateVendorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	vendor := &domain.Vendor{
		ID:            id,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
	}
	if req.Name != nil {
		vendor.Name = *req.Name
	}

	if err := h.svc.Update(ctx, user.ClinicID, vendor); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Vendor updated successfully"})
}

// DeleteVendor godoc
//
//	@Summary		Soft-delete a vendor
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Vendor ID"
//	@Success		200	{object}	response
//	@Router			/inventory/vendors/{id}/delete [patch]
func (h *VendorHandler) DeleteVendor(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	if err := h.svc.Delete(ctx, user.ClinicID, id); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Vendor deleted successfully"})
}

// ─── Stock Purchases ─────────────────────────────────────────────

type StockPurchaseHandler struct {
	svc port.StockPurchaseService
}

func NewStockPurchaseHandler(svc port.StockPurchaseService) *StockPurchaseHandler {
	return &StockPurchaseHandler{svc: svc}
}

// CreatePurchase godoc
//
//	@Summary		Record a stock purchase from a vendor
//	@Description	Creates the purchase and its line items, and increments quantity on hand for each item.
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateStockPurchaseReq	true	"Purchase details"
//	@Success		200		{object}	response{data=dto.StockPurchaseResponse}
//	@Router			/inventory/purchases [post]
func (h *StockPurchaseHandler) CreatePurchase(ctx *gin.Context) {
	var req dto.CreateStockPurchaseReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	purchase := &domain.StockPurchase{
		VendorID:     req.VendorID,
		InvoiceRefNo: req.InvoiceRefNo,
		PaidAmount:   req.PaidAmount,
		CreatedBy:    user.UserId,
	}

	items := make([]*domain.StockPurchaseItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, &domain.StockPurchaseItem{
			InventoryItemID: it.InventoryItemID,
			Quantity:        it.Quantity,
			UnitCost:        it.UnitCost,
		})
	}

	details, err := h.svc.CreatePurchase(ctx, user.ClinicID, purchase, items)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.StockPurchaseRes(details))
}

// ListPurchases godoc
//
//	@Summary		List stock purchases
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.StockPurchaseResponse}
//	@Router			/inventory/purchases [get]
func (h *StockPurchaseHandler) ListPurchases(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	purchases, err := h.svc.List(ctx, user.ClinicID, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.StockPurchaseResList(purchases))
}

// GetPurchaseByID godoc
//
//	@Summary		Get a stock purchase
//	@Tags			inventory
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Purchase ID"
//	@Success		200	{object}	response{data=dto.StockPurchaseResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/inventory/purchases/{id} [get]
func (h *StockPurchaseHandler) GetPurchaseByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	purchase, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.StockPurchaseRes(purchase))
}
