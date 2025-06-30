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
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		validationError(ctx, domain.ErrBadRequest)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	summary, err := h.svc.SalesDaily(ctx, user.ClinicID, date)
	if err != nil {
		handleError(ctx, err)
		return
	}

	res := dto.SalesSummaryRes(summary)
	headers, rows := res.TableRows()
	respondTable(ctx, "sales-daily-"+dateStr, "Daily Sales Report - "+dateStr, headers, rows, res)
}

// SalesMonthly godoc
//
//	@Summary		Monthly sales report
//	@Description	Admin-only. Add ?format=csv to download instead of receiving JSON.
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Param			year	query		int	false	"Defaults to current year"
//	@Param			month	query		int	false	"Defaults to current month"
//	@Success		200		{object}	response{data=dto.SalesSummaryResponse}
//	@Router			/reports/sales/monthly [get]
func (h *ReportsHandler) SalesMonthly(ctx *gin.Context) {
	now := time.Now()
	year, _ := strconv.Atoi(ctx.DefaultQuery("year", strconv.Itoa(now.Year())))
	month, _ := strconv.Atoi(ctx.DefaultQuery("month", strconv.Itoa(int(now.Month()))))

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	summary, err := h.svc.SalesMonthly(ctx, user.ClinicID, year, month)
	if err != nil {
		handleError(ctx, err)
		return
	}

	res := dto.SalesSummaryRes(summary)
	headers, rows := res.TableRows()
	respondTable(ctx, "sales-monthly", "Monthly Sales Report - "+res.Period, headers, rows, res)
}

// SalesRange godoc
//
//	@Summary		Sales report over a date range
//	@Description	Admin-only. Add ?format=csv to download instead of receiving JSON.
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Param			from	query		string	true	"YYYY-MM-DD"
//	@Param			to		query		string	true	"YYYY-MM-DD"
//	@Success		200		{object}	response{data=dto.SalesSummaryResponse}
//	@Router			/reports/sales/range [get]
func (h *ReportsHandler) SalesRange(ctx *gin.Context) {
	from, to, ok := parseDateRange(ctx)
	if !ok {
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	summary, err := h.svc.SalesRange(ctx, user.ClinicID, from, to)
	if err != nil {
		handleError(ctx, err)
		return
	}

	res := dto.SalesSummaryRes(summary)
	headers, rows := res.TableRows()
	respondTable(ctx, "sales-range", "Sales Report - "+res.Period, headers, rows, res)
}

// PatientDues godoc
//
//	@Summary		Outstanding patient balances
//	@Description	Admin-only. Add ?format=csv to download instead of receiving JSON.
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response{data=[]dto.PatientDueResponse}
