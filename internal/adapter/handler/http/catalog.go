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
	medicineSvc      port.MedicineService
	diagnosisSvc     port.DiagnosisCatalogService
	conditionSvc     port.HistoryConditionService
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

func (h *CatalogHandler) ListConditions(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	items, err := h.conditionSvc.List(ctx, limit, offset)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.HistoryConditionResList(items))
}

func (h *CatalogHandler) SearchConditions(ctx *gin.Context) {
	query := ctx.Query("query")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	items, err := h.conditionSvc.Search(ctx, query, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.HistoryConditionResList(items))
}

func (h *CatalogHandler) GetConditionByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	item, err := h.conditionSvc.GetByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.HistoryConditionRes(item))
}

func (h *CatalogHandler) UpdateCondition(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	var req dto.UpdateHistoryConditionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	item := &domain.HistoryCondition{
		ID:   id,
		Name: name,
	}

	err = h.conditionSvc.Update(ctx, item)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Condition updated successfully"})
}

func (h *CatalogHandler) DeleteCondition(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return
	}

	err = h.conditionSvc.Delete(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Condition deleted successfully"})
}
