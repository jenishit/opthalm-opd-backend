package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type CatalogHandler struct {
	medicineSvc  port.MedicineService
	diagnosisSvc port.DiagnosisCatalogService
	conditionSvc port.HistoryConditionService
}

func NewCatalogHandler(
	ms port.MedicineService,
	ds port.DiagnosisCatalogService,
	hs port.HistoryConditionService,
) *CatalogHandler {
	return &CatalogHandler{
		medicineSvc:  ms,
		diagnosisSvc: ds,
		conditionSvc: hs,
	}
}

// ─── Medicine ─────────────────────────────────────────────────────────────────

// ListMedicines godoc
//
//	@Summary		List medicines
//	@Description	Admin-only. Medicines are a global catalog, not clinic-scoped.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.MedicineResponse}
//	@Router			/admin/catalog/medicines [get]
func (h *CatalogHandler) ListMedicines(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	items, err := h.medicineSvc.List(ctx, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.MedicineResList(items))
}

// SearchMedicines godoc
//
//	@Summary		Search medicines
//	@Description	Admin-only.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.MedicineResponse}
//	@Router			/admin/catalog/medicines/search [get]
func (h *CatalogHandler) SearchMedicines(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	items, err := h.medicineSvc.Search(ctx, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.MedicineResList(items))
}

// GetMedicineByID godoc
//
//	@Summary		Get a medicine
//	@Description	Admin-only.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Medicine ID"
//	@Success		200	{object}	response{data=dto.MedicineResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/admin/catalog/medicines/{id} [get]
func (h *CatalogHandler) GetMedicineByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	item, err := h.medicineSvc.GetByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.MedicineRes(item))
}

// UpdateMedicine godoc
//
//	@Summary		Update a medicine
//	@Description	Admin-only.
//	@Tags			catalog
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Medicine ID"
//	@Param			request	body		dto.UpdateMedicineReq	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/admin/catalog/medicines/{id} [patch]
func (h *CatalogHandler) UpdateMedicine(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateMedicineReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	medicineName := ""
	if req.MedicineName != nil {
		medicineName = *req.MedicineName
	}

	item := &domain.Medicine{
		ID:           id,
		MedicineName: medicineName,
		BrandName:    req.BrandName,
		Strength:     req.Strength,
		Form:         req.Form,
	}

	err = h.medicineSvc.Update(ctx, item)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Medicine updated successfully"})
}

// DeleteMedicine godoc
//
//	@Summary		Soft-delete a medicine
//	@Description	Admin-only.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Medicine ID"
//	@Success		200	{object}	response
//	@Router			/admin/catalog/medicines/{id}/delete [patch]
func (h *CatalogHandler) DeleteMedicine(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	err = h.medicineSvc.Delete(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Medicine deleted successfully"})
}

// ─── Diagnosis Catalog ────────────────────────────────────────────────────────

// ListDiagnoses godoc
//
//	@Summary		List diagnoses
//	@Description	Admin-only. Diagnoses are a global catalog, not clinic-scoped.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.DiagnosisCatalogResponse}
//	@Router			/admin/catalog/diagnoses [get]
func (h *CatalogHandler) ListDiagnoses(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	items, err := h.diagnosisSvc.List(ctx, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.DiagnosisCatalogResList(items))
}

// SearchDiagnoses godoc
//
//	@Summary		Search diagnoses
//	@Description	Admin-only.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		string	false	"Search text"
//	@Param			limit	query		int		false	"Max results"	default(10)
//	@Success		200		{object}	response{data=[]dto.DiagnosisCatalogResponse}
//	@Router			/admin/catalog/diagnoses/search [get]
func (h *CatalogHandler) SearchDiagnoses(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	items, err := h.diagnosisSvc.Search(ctx, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.DiagnosisCatalogResList(items))
}

// GetDiagnosisByID godoc
//
//	@Summary		Get a diagnosis
//	@Description	Admin-only.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Diagnosis ID"
//	@Success		200	{object}	response{data=dto.DiagnosisCatalogResponse}
//	@Failure		404	{object}	errorResponse
//	@Router			/admin/catalog/diagnoses/{id} [get]
func (h *CatalogHandler) GetDiagnosisByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	item, err := h.diagnosisSvc.GetByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.DiagnosisCatalogRes(item))
}

// UpdateDiagnosis godoc
//
//	@Summary		Update a diagnosis
//	@Description	Admin-only.
//	@Tags			catalog
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Diagnosis ID"
//	@Param			request	body		dto.UpdateDiagnosisCatalogReq	true	"Fields to update"
//	@Success		200		{object}	response
//	@Router			/admin/catalog/diagnoses/{id} [patch]
func (h *CatalogHandler) UpdateDiagnosis(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateDiagnosisCatalogReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	item := &domain.DiagnosisCatalog{
		ID:        id,
		Name:      name,
		Icd10Code: req.Icd10Code,
	}

	err = h.diagnosisSvc.Update(ctx, item)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Diagnosis updated successfully"})
}

// DeleteDiagnosis godoc
//
//	@Summary		Soft-delete a diagnosis
//	@Description	Admin-only.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Diagnosis ID"
//	@Success		200	{object}	response
//	@Router			/admin/catalog/diagnoses/{id}/delete [patch]
func (h *CatalogHandler) DeleteDiagnosis(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	err = h.diagnosisSvc.Delete(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Diagnosis deleted successfully"})
}

// ─── History Condition ────────────────────────────────────────────────────────

// ListConditions godoc
//
//	@Summary		List history conditions
//	@Description	Admin-only. History conditions are a global catalog, not clinic-scoped.
//	@Tags			catalog
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Max results"	default(10)
//	@Param			offset	query		int	false	"Offset"		default(0)
//	@Success		200		{object}	response{data=[]dto.HistoryConditionResponse}
//	@Router			/admin/catalog/conditions [get]
func (h *CatalogHandler) ListConditions(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	items, err := h.conditionSvc.List(ctx, limit, offset)
	if err != nil {
		handleError(ctx, err)
