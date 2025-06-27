package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type ReportsHandler struct {
	svc port.ReportsService
}

func NewReportsHandler(svc port.ReportsService) *ReportsHandler {
	return &ReportsHandler{svc: svc}
}

// SalesDaily godoc
//
//	@Summary		Daily sales report
//	@Description	Admin-only. Add ?format=csv to download instead of receiving JSON.
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Param			date	query		string	false	"YYYY-MM-DD, defaults to today"
//	@Success		200		{object}	response{data=dto.SalesSummaryResponse}
//	@Router			/reports/sales/daily [get]
func (h *ReportsHandler) SalesDaily(ctx *gin.Context) {
	dateStr := ctx.DefaultQuery("date", time.Now().Format("2006-01-02"))
