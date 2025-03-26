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

