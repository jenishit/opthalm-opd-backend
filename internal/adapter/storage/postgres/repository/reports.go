package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ReportsRepository struct {
	DB *postgres.DB
}

func NewReportsRepository(db *postgres.DB) *ReportsRepository {
	return &ReportsRepository{DB: db}
}

func (r *ReportsRepository) SalesSummary(ctx context.Context, clinicID uuid.UUID, from, to time.Time) (*domain.SalesSummary, error) {
	query := `
		SELECT
			COUNT(*),
			COALESCE(SUM(subtotal), 0),
			COALESCE(SUM(discount_amount), 0),
			COALESCE(SUM(tax_amount), 0),
			COALESCE(SUM(total_amount), 0),
			COALESCE(SUM(paid_amount), 0),
			COALESCE(SUM(due_amount), 0)
		FROM invoices
		WHERE clinic_id = $1 AND status != 'cancelled' AND created_at >= $2 AND created_at < $3
	`

	s := &domain.SalesSummary{}
	err := r.DB.QueryRow(ctx, query, clinicID, from, to).Scan(
		&s.InvoiceCount, &s.Subtotal, &s.Discount, &s.Tax, &s.TotalSales, &s.TotalPaid, &s.TotalDue,
	)
	if err != nil {
		return nil, fmt.Errorf("ReportsRepo.SalesSummary: %w", err)
	}
	return s, nil
}

func (r *ReportsRepository) PatientDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.PatientDue, error) {
	query := `
		SELECT i.patient_id, p.full_name, p.phone, i.invoice_no, i.id, i.total_amount, i.paid_amount, i.due_amount, i.created_at
		FROM invoices i
		JOIN patients p ON p.id = i.patient_id
		WHERE i.clinic_id = $1 AND i.due_amount > 0 AND i.status != 'cancelled'
		ORDER BY i.created_at DESC
	`

	rows, err := r.DB.Query(ctx, query, clinicID)
	if err != nil {
		return nil, fmt.Errorf("ReportsRepo.PatientDues query: %w", err)
	}
	defer rows.Close()

	var dues []*domain.PatientDue
	for rows.Next() {
		var d domain.PatientDue
		if err := rows.Scan(&d.PatientID, &d.PatientName, &d.PatientPhone, &d.InvoiceNo, &d.InvoiceID, &d.TotalAmount, &d.PaidAmount, &d.DueAmount, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("ReportsRepo.PatientDues scan: %w", err)
		}
		dues = append(dues, &d)
	}
	return dues, rows.Err()
}

func (r *ReportsRepository) VendorDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.VendorDue, error) {
	query := `
		SELECT sp.vendor_id, v.name, sp.id, sp.total_amount, sp.paid_amount, sp.due_amount, sp.purchase_date
		FROM stock_purchases sp
		JOIN vendors v ON v.id = sp.vendor_id
		WHERE sp.clinic_id = $1 AND sp.due_amount > 0
		ORDER BY sp.purchase_date DESC
	`

	rows, err := r.DB.Query(ctx, query, clinicID)
	if err != nil {
		return nil, fmt.Errorf("ReportsRepo.VendorDues query: %w", err)
	}
	defer rows.Close()

	var dues []*domain.VendorDue
	for rows.Next() {
		var d domain.VendorDue
		if err := rows.Scan(&d.VendorID, &d.VendorName, &d.PurchaseID, &d.TotalAmount, &d.PaidAmount, &d.DueAmount, &d.PurchaseDate); err != nil {
			return nil, fmt.Errorf("ReportsRepo.VendorDues scan: %w", err)
		}
		dues = append(dues, &d)
	}
	return dues, rows.Err()
}

func (r *ReportsRepository) InventoryValuation(ctx context.Context, clinicID uuid.UUID) (*domain.InventoryValuation, error) {
	query := `
		SELECT
			COUNT(*),
			COALESCE(SUM(quantity_on_hand), 0),
