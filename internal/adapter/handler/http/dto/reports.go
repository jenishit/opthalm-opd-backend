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
		fmt.Sprintf("%.2f", s.TotalPaid), fmt.Sprintf("%.2f", s.TotalDue),
	}}
	return headers, rows
}

type PatientDueResponse struct {
	PatientName  string    `json:"patient_name"`
	PatientPhone string    `json:"patient_phone"`
	InvoiceNo    string    `json:"invoice_no"`
	TotalAmount  float64   `json:"total_amount"`
	PaidAmount   float64   `json:"paid_amount"`
	DueAmount    float64   `json:"due_amount"`
	CreatedAt    time.Time `json:"created_at"`
}

func PatientDueResList(ds []*domain.PatientDue) []*PatientDueResponse {
	res := make([]*PatientDueResponse, 0, len(ds))
	for _, d := range ds {
		res = append(res, &PatientDueResponse{
			PatientName: d.PatientName, PatientPhone: d.PatientPhone, InvoiceNo: d.InvoiceNo,
			TotalAmount: d.TotalAmount, PaidAmount: d.PaidAmount, DueAmount: d.DueAmount, CreatedAt: d.CreatedAt,
		})
	}
	return res
}

func PatientDueTableRows(ds []*PatientDueResponse) ([]string, [][]string) {
	headers := []string{"Patient", "Phone", "Invoice No", "Total", "Paid", "Due", "Date"}
	rows := make([][]string, 0, len(ds))
	for _, d := range ds {
		rows = append(rows, []string{
			d.PatientName, d.PatientPhone, d.InvoiceNo,
			fmt.Sprintf("%.2f", d.TotalAmount), fmt.Sprintf("%.2f", d.PaidAmount), fmt.Sprintf("%.2f", d.DueAmount),
			d.CreatedAt.Format("2006-01-02"),
		})
	}
	return headers, rows
}

type VendorDueResponse struct {
	VendorName   string    `json:"vendor_name"`
	TotalAmount  float64   `json:"total_amount"`
	PaidAmount   float64   `json:"paid_amount"`
	DueAmount    float64   `json:"due_amount"`
	PurchaseDate time.Time `json:"purchase_date"`
}

func VendorDueResList(ds []*domain.VendorDue) []*VendorDueResponse {
	res := make([]*VendorDueResponse, 0, len(ds))
	for _, d := range ds {
		res = append(res, &VendorDueResponse{
			VendorName: d.VendorName, TotalAmount: d.TotalAmount, PaidAmount: d.PaidAmount,
			DueAmount: d.DueAmount, PurchaseDate: d.PurchaseDate,
		})
	}
	return res
}

func VendorDueTableRows(ds []*VendorDueResponse) ([]string, [][]string) {
	headers := []string{"Vendor", "Total", "Paid", "Due", "Purchase Date"}
	rows := make([][]string, 0, len(ds))
	for _, d := range ds {
		rows = append(rows, []string{
