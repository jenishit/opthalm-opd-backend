package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type LabJobRepository struct {
	DB *postgres.DB
}

func NewLabJobRepository(db *postgres.DB) *LabJobRepository {
	return &LabJobRepository{DB: db}
}

func (r *LabJobRepository) Create(ctx context.Context, clinicID uuid.UUID, job *domain.LabJob) (*domain.LabJobDetails, error) {
	if job.Status == "" {
		job.Status = domain.LabJobInFitting
	}

	err := r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		query, args, err := sq.Insert("lab_jobs").
			Columns(
				"clinic_id", "invoice_id", "invoice_item_id", "patient_id", "vendor_id", "job_type", "status",
				"expected_delivery_date", "advance_payment", "notes", "created_by", "updated_by",
