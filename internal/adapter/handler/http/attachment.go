package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

// maxAttachmentUploadSize caps a single clinical image upload (fundus
// photos, OCT/visual-field exports, biometry printouts) — generous enough
// for a camera-resolution photo or a multi-page OCT PDF, without letting
// one request exhaust the server's memory.
const maxAttachmentUploadSize = 25 << 20 // 25 MiB

var allowedAttachmentTypes = map[string]bool{
	"image/png":       true,
	"image/jpeg":      true,
	"application/pdf": true,
}

type AttachmentHandler struct {
	svc      port.AttachmentService
	visitSvc port.VisitsService
}

func NewAttachmentHandler(svc port.AttachmentService, visitSvc port.VisitsService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc, visitSvc: visitSvc}
}

// Upload godoc
//
//	@Summary		Upload a clinical image for a visit
//	@Description	Images pulled off diagnostic devices (fundus camera, OCT, visual field analyzer, biometry) — download from the device, then upload here. PNG, JPEG, or PDF, up to 25MB.
//	@Tags			clinical
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		string	true	"Visit ID"
//	@Param			file		formData	file	true	"Image file (PNG/JPEG/PDF)"
//	@Param			category	formData	string	true	"fundus_photo | oct | visual_field | biometry | topography | other"
//	@Param			eye			formData	string	false	"OD | OS | OU | none"
//	@Param			notes		formData	string	false	"Notes"
//	@Success		200			{object}	response{data=dto.AttachmentResponse}
//	@Router			/visit/{id}/attachments [post]
func (h *AttachmentHandler) Upload(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.UploadAttachmentReq
	if err := ctx.ShouldBind(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	visit, err := h.visitSvc.GetVisitByVisitID(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxAttachmentUploadSize)

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		validationError(ctx, fmt.Errorf("file is required: %w", err))
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedAttachmentTypes[contentType] {
		validationError(ctx, domain.ErrBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		handleError(ctx, fmt.Errorf("opening upload: %w", err))
		return
	}
	defer file.Close()

	eye := req.Eye
	if eye == "" {
		eye = "none"
	}

	attachment, url, err := h.svc.UploadAttachment(
		ctx,
		user.ClinicID,
		visitID,
		visit.PatientID,
		user.UserId,
		domain.AttachmentCategory(req.Category),
		eye,
		fileHeader.Filename,
		contentType,
		fileHeader.Size,
		file,
		req.Notes,
	)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.AttachmentRes(attachment, url))
}

// ListAttachments godoc
//
//	@Summary		List a visit's clinical images
//	@Tags			clinical
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Visit ID"
//	@Success		200	{object}	response{data=[]dto.AttachmentResponse}
//	@Router			/visit/{id}/attachments [get]
func (h *AttachmentHandler) ListAttachments(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	attachments, err := h.svc.ListAttachments(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := make([]*dto.AttachmentResponse, 0, len(attachments))
	for _, a := range attachments {
		rsp = append(rsp, dto.AttachmentRes(a.ClinicalAttachment, a.URL))
	}

	handleSuccess(ctx, rsp)
}

// DeleteAttachment godoc
//
//	@Summary		Delete a clinical image
//	@Tags			clinical
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		string	true	"Visit ID"
//	@Param			attachmentId	path		string	true	"Attachment ID"
//	@Success		200				{object}	response
//	@Router			/visit/{id}/attachments/{attachmentId} [delete]
func (h *AttachmentHandler) DeleteAttachment(ctx *gin.Context) {
	attachmentID, err := uuid.Parse(ctx.Param("attachmentId"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	if err := h.svc.DeleteAttachment(ctx, user.ClinicID, attachmentID); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Attachment deleted successfully"})
}
