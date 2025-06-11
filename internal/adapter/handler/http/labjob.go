package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type LabJobHandler struct {
	svc port.LabJobService
}

func NewLabJobHandler(svc port.LabJobService) *LabJobHandler {
	return &LabJobHandler{svc: svc}
}

// CreateLabJob godoc
//
//	@Summary		Create a lab job
//	@Description	Requires ROLE_ADMIN, ROLE_LAB, or ROLE_BILLING.
//	@Tags			lab-jobs
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateLabJobReq	true	"Lab job details"
//	@Success		200		{object}	response{data=dto.LabJobResponse}
//	@Router			/lab-jobs [post]
func (h *LabJobHandler) CreateLabJob(ctx *gin.Context) {
	var req dto.CreateLabJobReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	var expectedDeliveryDate *time.Time
	if req.ExpectedDeliveryDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.ExpectedDeliveryDate)
		if err != nil {
			validationError(ctx, domain.ErrBadRequest)
			return
		}
		expectedDeliveryDate = &parsed
	}

	job := &domain.LabJob{
		InvoiceID:            req.InvoiceID,
		InvoiceItemID:        req.InvoiceItemID,
		PatientID:            req.PatientID,
		VendorID:             req.VendorID,
		JobType:              req.JobType,
		ExpectedDeliveryDate: expectedDeliveryDate,
		AdvancePayment:       req.AdvancePayment,
		Notes:                req.Notes,
		CreatedBy:            user.UserId,
		UpdatedBy:            user.UserId,
	}

	created, err := h.svc.Create(ctx, user.ClinicID, job)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.LabJobRes(created))
}

// ListLabJobs godoc
//
//	@Summary		List lab jobs
//	@Tags			lab-jobs
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.LabJobResponse}
//	@Router			/lab-jobs [get]
func (h *LabJobHandler) ListLabJobs(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	jobs, err := h.svc.List(ctx, user.ClinicID, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.LabJobResList(jobs))
}

// GetLabJobByID godoc
//
//	@Summary		Get a lab job
//	@Tags			lab-jobs
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Lab job ID"
//	@Success		200	{object}	response{data=dto.LabJobResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/lab-jobs/{id} [get]
func (h *LabJobHandler) GetLabJobByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	job, err := h.svc.GetByID(ctx, user.ClinicID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.LabJobRes(job))
}
