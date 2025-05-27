package dto

import (
	"fmt"
	"time"

	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type SalesSummaryResponse struct {
	Period       string  `json:"period"`
	InvoiceCount int     `json:"invoice_count"`
	Subtotal     float64 `json:"subtotal"`
	Discount     float64 `json:"discount"`
	Tax          float64 `json:"tax"`
	TotalSales   float64 `json:"total_sales"`
	TotalPaid    float64 `json:"total_paid"`
	TotalDue     float64 `json:"total_due"`
}

func SalesSummaryRes(s *domain.SalesSummary) *SalesSummaryResponse {
	return &SalesSummaryResponse{
		Period: s.Period, InvoiceCount: s.InvoiceCount, Subtotal: s.Subtotal,
		Discount: s.Discount, Tax: s.Tax, TotalSales: s.TotalSales, TotalPaid: s.TotalPaid, TotalDue: s.TotalDue,
	}
}

func (s *SalesSummaryResponse) TableRows() ([]string, [][]string) {
	headers := []string{"Period", "Invoices", "Subtotal", "Discount", "Tax", "Total", "Paid", "Due"}
	rows := [][]string{{
		s.Period, fmt.Sprintf("%d", s.InvoiceCount), fmt.Sprintf("%.2f", s.Subtotal),
		fmt.Sprintf("%.2f", s.Discount), fmt.Sprintf("%.2f", s.Tax), fmt.Sprintf("%.2f", s.TotalSales),
