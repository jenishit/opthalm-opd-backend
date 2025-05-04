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
