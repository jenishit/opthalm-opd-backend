package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

// ClinicConfigHandler serves the Settings pages that are a single record
// (or a short list) per clinic: locations, notification preferences,
// patient-facing payment/tax config, barcode & POS config, data & privacy
// config plus its audit log, and third-party integrations.
type ClinicConfigHandler struct {
	locationSvc     port.ClinicLocationService
	notificationSvc port.NotificationPrefsService
	paymentSvc      port.PaymentConfigService
	posSvc          port.PosConfigService
	privacySvc      port.PrivacyConfigService
	auditLogSvc     port.AuditLogService
	integrationSvc  port.IntegrationService
	subscriptionSvc port.SubscriptionService
	uploadDir       string
}

func NewClinicConfigHandler(
	locationSvc port.ClinicLocationService,
	notificationSvc port.NotificationPrefsService,
	paymentSvc port.PaymentConfigService,
	posSvc port.PosConfigService,
	privacySvc port.PrivacyConfigService,
	auditLogSvc port.AuditLogService,
	integrationSvc port.IntegrationService,
	subscriptionSvc port.SubscriptionService,
	uploadDir string,
) *ClinicConfigHandler {
	return &ClinicConfigHandler{
		locationSvc:     locationSvc,
		notificationSvc: notificationSvc,
		paymentSvc:      paymentSvc,
		posSvc:          posSvc,
		privacySvc:      privacySvc,
		auditLogSvc:     auditLogSvc,
		integrationSvc:  integrationSvc,
		subscriptionSvc: subscriptionSvc,
		uploadDir:       uploadDir,
	}
}

// ─── Billing & subscription ─────────────────────────────────────────────────

// GetBillingSubscription godoc
//
//	@Summary		Get the clinic's own subscription, plan limits, and live usage
//	@Description	Unlike /platform/subscriptions/{clinicId}, this is scoped to the caller's own clinic — any admin can read it, not just a platform operator.
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=dto.BillingSubscriptionResponse}
//	@Router			/admin/clinic/billing/subscription [get]
func (h *ClinicConfigHandler) GetBillingSubscription(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	sub, err := h.subscriptionSvc.GetByClinicID(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	usage, err := h.subscriptionSvc.GetUsageMetrics(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.BillingSubscriptionRes(sub, usage))
}

// ─── Locations & rooms ──────────────────────────────────────────────────────

// ListLocations godoc
//
//	@Summary		List the clinic's locations
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=[]dto.ClinicLocationResponse}
//	@Router			/admin/clinic/locations [get]
func (h *ClinicConfigHandler) ListLocations(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	locations, err := h.locationSvc.ListLocations(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.ClinicLocationResList(locations))
}

// CreateLocation godoc
//
//	@Summary		Add a clinic location
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateClinicLocationReq	true	"Location"
//	@Success		200		{object}	response{data=dto.ClinicLocationResponse}
//	@Router			/admin/clinic/locations [post]
func (h *ClinicConfigHandler) CreateLocation(ctx *gin.Context) {
	var req dto.CreateClinicLocationReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	location := &domain.ClinicLocation{
		Name:      req.Name,
		Address:   req.Address,
		ExamRooms: req.ExamRooms,
		Status:    req.Status,
		CreatedBy: user.UserId,
	}

	saved, err := h.locationSvc.CreateLocation(ctx, user.ClinicID, location)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.ClinicLocationRes(saved))
}

// ─── Notification preferences ──────────────────────────────────────────────

// GetNotificationPrefs godoc
//
//	@Summary		Get notification preferences
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=dto.NotificationPrefsResponse}
//	@Router			/admin/clinic/notifications [get]
func (h *ClinicConfigHandler) GetNotificationPrefs(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	prefs, err := h.notificationSvc.GetNotificationPrefs(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.NotificationPrefsRes(prefs))
}

// UpdateNotificationPrefs godoc
//
//	@Summary		Update notification preferences
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpsertNotificationPrefsReq	true	"Preferences"
//	@Success		200		{object}	response{data=dto.NotificationPrefsResponse}
//	@Router			/admin/clinic/notifications [put]
func (h *ClinicConfigHandler) UpdateNotificationPrefs(ctx *gin.Context) {
	var req dto.UpsertNotificationPrefsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	prefs := &domain.NotificationPrefs{
		ClinicID:                user.ClinicID,
		EmailNotifications:      req.EmailNotifications,
		SMSAppointmentReminders: req.SMSAppointmentReminders,
		LowStockAlerts:          req.LowStockAlerts,
		NewOrderNotifications:   req.NewOrderNotifications,
		WeeklyDigest:            req.WeeklyDigest,
	}

	saved, err := h.notificationSvc.UpsertNotificationPrefs(ctx, prefs)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.NotificationPrefsRes(saved))
}

// ─── Payment methods & tax config ───────────────────────────────────────────

// GetPaymentConfig godoc
//
//	@Summary		Get the clinic's patient-facing payment & tax config
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=dto.PaymentConfigResponse}
//	@Router			/admin/clinic/payment-config [get]
func (h *ClinicConfigHandler) GetPaymentConfig(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	cfg, err := h.paymentSvc.GetPaymentConfig(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PaymentConfigRes(cfg))
}

// UpdatePaymentConfig godoc
//
//	@Summary		Update the clinic's patient-facing payment & tax config
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpsertPaymentConfigReq	true	"Payment config"
//	@Success		200		{object}	response{data=dto.PaymentConfigResponse}
//	@Router			/admin/clinic/payment-config [put]
func (h *ClinicConfigHandler) UpdatePaymentConfig(ctx *gin.Context) {
	var req dto.UpsertPaymentConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	cfg := &domain.PaymentConfig{
		ClinicID:  user.ClinicID,
		TaxLabel:  req.TaxLabel,
		TaxRate:   req.TaxRate,
		VATNumber: req.VATNumber,
		PANNumber: req.PANNumber,
		Methods: domain.PaymentMethods{
			Cash:         req.Methods.Cash,
			Card:         req.Methods.Card,
			BankTransfer: req.Methods.BankTransfer,
			QRWalletFlag: req.Methods.QRWallet,
		},
	}

	if req.BankDetail != nil {
		cfg.BankDetail = &domain.BankDetails{
			BankName:      req.BankDetail.BankName,
			AccountHolder: req.BankDetail.AccountHolder,
			AccountNumber: req.BankDetail.AccountNumber,
			Branch:        req.BankDetail.Branch,
			SwiftCode:     req.BankDetail.SwiftCode,
		}
	}

	if req.QRWallet != nil {
		cfg.QRWallet = &domain.QRWallet{
			Provider:   req.QRWallet.Provider,
			MerchantID: req.QRWallet.MerchantID,
			QRAssetURL: req.QRWallet.QRAssetURL,
		}
	}

	saved, err := h.paymentSvc.UpsertPaymentConfig(ctx, cfg)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PaymentConfigRes(saved))
}

// maxQRUploadSize caps the multipart body gin will read into memory before
// spilling the rest to a temp file; a static QR image has no business being
// larger than this.
const maxQRUploadSize = 5 << 20 // 5 MiB

var allowedQRImageTypes = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
}

// UploadQRCode godoc
//
//	@Summary		Upload the clinic's static payment QR code image
//	@Tags			clinic
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			file	formData	file	true	"QR code image (PNG or JPEG)"
//	@Success		200		{object}	response{data=dto.PaymentConfigResponse}
//	@Router			/admin/clinic/payment-config/qr-code [post]
func (h *ClinicConfigHandler) UploadQRCode(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxQRUploadSize)

	file, err := ctx.FormFile("file")
	if err != nil {
		validationError(ctx, fmt.Errorf("file is required: %w", err))
		return
	}

	ext, ok := allowedQRImageTypes[file.Header.Get("Content-Type")]
	if !ok {
		validationError(ctx, domain.ErrBadRequest)
		return
	}

	dir := filepath.Join(h.uploadDir, "clinics", user.ClinicID.String())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		handleError(ctx, fmt.Errorf("creating upload dir: %w", err))
		return
	}

	dst := filepath.Join(dir, "qr-code"+ext)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		handleError(ctx, fmt.Errorf("saving upload: %w", err))
		return
	}

	url := "/uploads/clinics/" + user.ClinicID.String() + "/qr-code" + ext

	saved, err := h.paymentSvc.SetQRAssetURL(ctx, user.ClinicID, url)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PaymentConfigRes(saved))
}

// ─── Barcode & POS config ───────────────────────────────────────────────────

// GetPosConfig godoc
//
//	@Summary		Get the clinic's barcode & POS config
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=dto.PosConfigResponse}
//	@Router			/admin/clinic/pos-config [get]
func (h *ClinicConfigHandler) GetPosConfig(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	cfg, err := h.posSvc.GetPosConfig(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PosConfigRes(cfg))
}

// UpdatePosConfig godoc
//
//	@Summary		Update the clinic's barcode & POS config
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpsertPosConfigReq	true	"POS config"
//	@Success		200		{object}	response{data=dto.PosConfigResponse}
//	@Router			/admin/clinic/pos-config [put]
func (h *ClinicConfigHandler) UpdatePosConfig(ctx *gin.Context) {
	var req dto.UpsertPosConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	cfg := &domain.PosConfig{
		ClinicID:               user.ClinicID,
		BarcodeSymbology:       req.BarcodeSymbology,
		AutoGenerateBarcode:    req.AutoGenerateBarcode,
		BarcodeScanningEnabled: req.BarcodeScanningEnabled,
		AutoPrintReceipt:       req.AutoPrintReceipt,
		RequirePinForRefunds:   req.RequirePinForRefunds,
		ReceiptPrinter:         req.ReceiptPrinter,
		TillFloatAmount:        req.TillFloatAmount,
	}

	saved, err := h.posSvc.UpsertPosConfig(ctx, cfg)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PosConfigRes(saved))
}

// ─── Data & privacy config ──────────────────────────────────────────────────

// GetPrivacyConfig godoc
//
//	@Summary		Get the clinic's data & privacy config
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=dto.PrivacyConfigResponse}
//	@Router			/admin/clinic/privacy-config [get]
func (h *ClinicConfigHandler) GetPrivacyConfig(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	cfg, err := h.privacySvc.GetPrivacyConfig(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PrivacyConfigRes(cfg))
}

// UpdatePrivacyConfig godoc
//
//	@Summary		Update the clinic's data & privacy config
//	@Tags			clinic
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpsertPrivacyConfigReq	true	"Privacy config"
//	@Success		200		{object}	response{data=dto.PrivacyConfigResponse}
//	@Router			/admin/clinic/privacy-config [put]
func (h *ClinicConfigHandler) UpdatePrivacyConfig(ctx *gin.Context) {
	var req dto.UpsertPrivacyConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	cfg := &domain.PrivacyConfig{
		ClinicID:                     user.ClinicID,
		RetentionYears:               req.RetentionYears,
		AutoAnonymizeAfterRetention:  req.AutoAnonymizeAfterRetention,
		AllowPatientDataExport:       req.AllowPatientDataExport,
		AllowPatientDeletionRequests: req.AllowPatientDeletionRequests,
		Require2FA:                   req.Require2FA,
		SessionTimeoutMinutes:        req.SessionTimeoutMinutes,
		LoginAlerts:                  req.LoginAlerts,
	}

	saved, err := h.privacySvc.UpsertPrivacyConfig(ctx, cfg)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.PrivacyConfigRes(saved))
}

// ListAuditLog godoc
//
//	@Summary		List the clinic's audit log, newest first
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(50)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.AuditLogEntryResponse}
//	@Router			/admin/clinic/audit-log [get]
func (h *ClinicConfigHandler) ListAuditLog(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	entries, err := h.auditLogSvc.ListAuditLog(ctx, user.ClinicID, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.AuditLogEntryResList(entries))
}

// ─── Integrations ────────────────────────────────────────────────────────────

// ListIntegrations godoc
//
//	@Summary		List the clinic's third-party integrations
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=[]dto.IntegrationResponse}
//	@Router			/admin/clinic/integrations [get]
func (h *ClinicConfigHandler) ListIntegrations(ctx *gin.Context) {
	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	integrations, err := h.integrationSvc.ListIntegrations(ctx, user.ClinicID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.IntegrationResList(integrations))
}

// ConnectIntegration godoc
//
//	@Summary		Start connecting a third-party integration
//	@Description	Stub: marks the integration connected immediately. A real OAuth redirect+callback flow is not wired up yet.
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Param			provider	path		string	true	"Provider"	Enums(xero, twilio, google_calendar, essilor_connect)
//	@Success		200			{object}	response{data=dto.IntegrationResponse}
//	@Router			/admin/clinic/integrations/{provider}/connect [post]
func (h *ClinicConfigHandler) ConnectIntegration(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if !domain.IsKnownIntegrationProvider(provider) {
		handleError(ctx, domain.ErrBadRequest)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	integration, err := h.integrationSvc.ConnectIntegration(ctx, user.ClinicID, provider)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.IntegrationRes(integration))
}

// DisconnectIntegration godoc
//
//	@Summary		Disconnect a third-party integration
//	@Tags			clinic
//	@Produce		json
//	@Security		BearerAuth
//	@Param			provider	path		string	true	"Provider"	Enums(xero, twilio, google_calendar, essilor_connect)
//	@Success		200			{object}	response
//	@Router			/admin/clinic/integrations/{provider} [delete]
func (h *ClinicConfigHandler) DisconnectIntegration(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if !domain.IsKnownIntegrationProvider(provider) {
		handleError(ctx, domain.ErrBadRequest)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	if err := h.integrationSvc.DisconnectIntegration(ctx, user.ClinicID, provider); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Integration disconnected successfully"})
}
