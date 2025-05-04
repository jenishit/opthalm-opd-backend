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
