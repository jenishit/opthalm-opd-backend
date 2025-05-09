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
