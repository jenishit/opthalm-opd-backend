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
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InvoiceResList(invoices))
}

// GetInvoiceByID godoc
//
//	@Summary		Get an invoice
//	@Tags			billing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Invoice ID"
//	@Success		200	{object}	response{data=dto.InvoiceResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/billing/invoice/{id} [get]
func (h *InvoiceHandler) GetInvoiceByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoice, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.InvoiceRes(invoice))
}

// UpdateInvoiceStatus godoc
//
//	@Summary		Update invoice status
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Invoice ID"
//	@Param			request	body		dto.UpdateInvoiceStatusReq	true	"New status"
//	@Success		200		{object}	response
//	@Failure		409		{object}	errorResponse	"invoice already paid"
//	@Router			/billing/invoice/{id}/status [patch]
func (h *InvoiceHandler) UpdateInvoiceStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateInvoiceStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	if err := h.svc.UpdateStatus(ctx, user.ClinicID, id, domain.InvoiceStatus(req.Status), user.UserId); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Invoice status updated successfully"})
}

// RecordPayment godoc
//
//	@Summary		Record a payment against an invoice
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Invoice ID"
//	@Param			request	body		dto.RecordPaymentReq	true	"Payment details"
//	@Success		200		{object}	response{data=dto.PaymentResponse}
//	@Failure		400		{object}	errorResponse	"payment exceeds due amount"
//	@Router			/billing/invoice/{id}/payments [post]
func (h *InvoiceHandler) RecordPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.RecordPaymentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	payment := &domain.Payment{
		InvoiceID:   id,
		Amount:      req.Amount,
		Method:      domain.PaymentMethod(req.Method),
		ReferenceNo: req.ReferenceNo,
		CreatedBy:   user.UserId,
	}

	p, err := h.svc.RecordPayment(ctx, user.ClinicID, payment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PaymentRes(p))
}

// GetInvoicePDF godoc
//
//	@Summary		Download an invoice as PDF
//	@Tags			billing
//	@Produce		application/pdf
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Invoice ID"
//	@Success		200	{file}	binary
//	@Router			/billing/invoice/{id}/pdf [get]
func (h *InvoiceHandler) GetInvoicePDF(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoice, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Branding is best-effort: an unconfigured clinic still gets a usable invoice.
	var clinic *domain.ClinicSettings
	if c, err := h.clinicSvc.GetClinicByClinicID(ctx, user.ClinicID); err == nil {
		clinic = c
	}

	buf, err := renderInvoicePDF(invoice, clinic)
	if err != nil {
		handleError(ctx, domain.ErrInternal)
		return
	}

	ctx.Data(200, "application/pdf", buf.Bytes())
}

// GetInvoiceQR godoc
//
//	@Summary		Get an invoice summary as a QR code
//	@Tags			billing
//	@Produce		image/png
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Invoice ID"
//	@Success		200	{file}	binary
//	@Router			/billing/invoice/{id}/qr [get]
func (h *InvoiceHandler) GetInvoiceQR(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoice, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	content := fmt.Sprintf("Invoice %s\nTotal: %.2f\nPaid: %.2f\nDue: %.2f",
		invoice.InvoiceNo, invoice.TotalAmount, invoice.PaidAmount, invoice.DueAmount)

	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		handleError(ctx, domain.ErrInternal)
		return
	}

	ctx.Data(200, "image/png", png)
}

// GetInvoiceWhatsAppLink godoc
//
//	@Summary		Get a wa.me link to share an invoice summary
//	@Tags			billing
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Invoice ID"
//	@Success		200	{object}	response{data=object{link=string}}
//	@Router			/billing/invoice/{id}/whatsapp-link [get]
func (h *InvoiceHandler) GetInvoiceWhatsAppLink(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	invoice, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	text := fmt.Sprintf(
		"Invoice %s\nTotal: %.2f\nPaid: %.2f\nDue: %.2f\nThank you for your visit!",
		invoice.InvoiceNo, invoice.TotalAmount, invoice.PaidAmount, invoice.DueAmount,
	)

	// PatientPhone is stored as entered at registration; the frontend/operator
	// is responsible for making sure it's in a WhatsApp-resolvable format
	// (wa.me tolerates a leading '+' and digits, and falls back to a contact
	// picker if the number can't be resolved).
	link := fmt.Sprintf("https://wa.me/%s?text=%s", url.QueryEscape(invoice.PatientPhone), url.QueryEscape(text))

	handleSuccess(ctx, gin.H{"link": link})
}

func renderInvoicePDF(invoice *domain.InvoiceDetails, clinic *domain.ClinicSettings) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	clinicName := "Invoice"
	if clinic != nil && clinic.ClinicName != "" {
		clinicName = clinic.ClinicName
	}
	pdf.Cell(0, 10, clinicName)
	pdf.Ln(8)

	if clinic != nil {
		pdf.SetFont("Arial", "", 10)
		if clinic.Address != nil {
			pdf.Cell(0, 6, *clinic.Address)
			pdf.Ln(5)
		}
		if clinic.Phone != nil {
			pdf.Cell(0, 6, "Phone: "+*clinic.Phone)
			pdf.Ln(5)
		}
	}
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Invoice "+invoice.InvoiceNo)
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "Date: "+invoice.CreatedAt.Format("2006-01-02"))
	pdf.Ln(5)
	pdf.Cell(0, 6, "Patient: "+invoice.PatientName+"  ("+invoice.PatientPhone+")")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(90, 7, "Description", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 7, "Qty", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, "Unit Price", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, "Line Total", "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, item := range invoice.Items {
		pdf.CellFormat(90, 7, item.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 7, fmt.Sprintf("%d", item.Quantity), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", item.LineTotal), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 10)
	totalsLine := func(label string, value float64) {
		pdf.CellFormat(140, 7, "", "", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, label, "", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", value), "", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	totalsLine("Subtotal", invoice.Subtotal)
	totalsLine("Discount", invoice.DiscountAmount)
	totalsLine("Tax", invoice.TaxAmount)
	totalsLine("Total", invoice.TotalAmount)
	totalsLine("Paid", invoice.PaidAmount)
	totalsLine("Due", invoice.DueAmount)

	if clinic != nil && clinic.ReportFooter != nil {
		pdf.Ln(10)
		pdf.SetFont("Arial", "I", 8)
		pdf.MultiCell(0, 5, *clinic.ReportFooter, "", "L", false)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}
