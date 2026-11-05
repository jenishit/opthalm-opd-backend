package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

// maxSignatureUploadSize is small — a signature is a flat-color PNG
// exported from a canvas, not a photo, so it's nowhere near a clinical
// attachment's size.
const maxSignatureUploadSize = 2 << 20 // 2 MiB

var allowedSignatureTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
}

const defaultSignedByRelation = "self"

type ConsentHandler struct {
	svc      port.ConsentService
	visitSvc port.VisitsService
}

func NewConsentHandler(svc port.ConsentService, visitSvc port.VisitsService) *ConsentHandler {
	return &ConsentHandler{svc: svc, visitSvc: visitSvc}
}

// Create godoc
//
//	@Summary		Record a signed consent for a visit
//	@Description	The signature is captured on a signature pad or tablet canvas client-side, exported as a PNG/JPEG image, and uploaded here alongside who signed it.
//	@Tags			clinical
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id					path		string	true	"Visit ID"
//	@Param			signature			formData	file	true	"Signature image (PNG/JPEG)"
//	@Param			consent_type		formData	string	true	"exam | procedure | surgery | general"
//	@Param			signed_by_name		formData	string	true	"Name of the person who signed"
//	@Param			signed_by_relation	formData	string	false	"self | guardian | other (default self)"
//	@Param			witnessed_by		formData	string	false	"Staff user ID who witnessed the signature"
//	@Param			notes				formData	string	false	"Notes"
//	@Success		200					{object}	response{data=dto.ConsentResponse}
//	@Router			/visit/{id}/consents [post]
func (h *ConsentHandler) Create(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.RecordConsentReq
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

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSignatureUploadSize)

	fileHeader, err := ctx.FormFile("signature")
	if err != nil {
		validationError(ctx, fmt.Errorf("signature file is required: %w", err))
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedSignatureTypes[contentType] {
		validationError(ctx, domain.ErrBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		handleError(ctx, fmt.Errorf("opening signature upload: %w", err))
		return
	}
	defer file.Close()

	relation := req.SignedByRelation
	if relation == "" {
		relation = defaultSignedByRelation
	}

	consent, url, err := h.svc.RecordConsent(
		ctx,
		user.ClinicID,
		visitID,
		visit.PatientID,
		user.UserId,
		domain.ConsentType(req.ConsentType),
		req.SignedByName,
		domain.SignedByRelation(relation),
		req.WitnessedBy,
		req.Notes,
		contentType,
		fileHeader.Size,
		file,
	)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.ConsentRes(consent, url))
}

// List godoc
//
//	@Summary		List a visit's signed consents
//	@Tags			clinical
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Visit ID"
//	@Success		200	{object}	response{data=[]dto.ConsentResponse}
//	@Router			/visit/{id}/consents [get]
func (h *ConsentHandler) List(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	consents, err := h.svc.ListConsents(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := make([]*dto.ConsentResponse, 0, len(consents))
	for _, c := range consents {
		rsp = append(rsp, dto.ConsentRes(c.Consent, c.SignatureURL))
	}

	handleSuccess(ctx, rsp)
}
