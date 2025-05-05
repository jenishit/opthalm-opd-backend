package http

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

type InvoiceHandler struct {
	svc       port.InvoiceService
	clinicSvc port.ClinicService
}

func NewInvoiceHandler(svc port.InvoiceService, clinicSvc port.ClinicService) *InvoiceHandler {
	return &InvoiceHandler{svc: svc, clinicSvc: clinicSvc}
}

// CreateInvoice godoc
//
//	@Summary		Create an invoice
//	@Description	Requires ROLE_ADMIN or ROLE_BILLING.
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateInvoiceReq	true	"Invoice details"
//	@Success		200		{object}	response{data=dto.InvoiceResponse}
//	@Failure		409		{object}	errorResponse	"insufficient stock"
//	@Router			/billing/invoice [post]
func (h *InvoiceHandler) CreateInvoice(ctx *gin.Context) {
	var req dto.CreateInvoiceReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoice := &domain.Invoice{
		PatientID:      req.PatientID,
		VisitID:        req.VisitID,
		DiscountAmount: req.DiscountAmount,
		TaxAmount:      req.TaxAmount,
		CreatedBy:      user.UserId,
		UpdatedBy:      user.UserId,
	}

	items := make([]*domain.InvoiceItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, &domain.InvoiceItem{
			BundleID:        it.BundleID,
			ItemType:        domain.ItemType(it.ItemType),
			Description:     it.Description,
			InventoryItemID: it.InventoryItemID,
			Quantity:        it.Quantity,
			UnitPrice:       it.UnitPrice,
			DiscountAmount:  it.DiscountAmount,
		})
	}

	invoiceDetails, err := h.svc.CreateInvoice(ctx, user.ClinicID, invoice, items)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InvoiceRes(invoiceDetails))
}

// ListInvoices godoc
//
//	@Summary		List invoices
//	@Tags			billing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.InvoiceResponse}
//	@Router			/billing/invoice [get]
func (h *InvoiceHandler) ListInvoices(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoices, err := h.svc.List(ctx, user.ClinicID, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InvoiceResList(invoices))
}

// SearchInvoices godoc
//
//	@Summary		Search invoices
//	@Tags			billing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.InvoiceResponse}
//	@Router			/billing/invoice/search [get]
func (h *InvoiceHandler) SearchInvoices(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoices, err := h.svc.Search(ctx, user.ClinicID, query, limit)
