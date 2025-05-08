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

