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
			).
			Values(
				clinicID, nullUUIDPtr(job.InvoiceID), nullUUIDPtr(job.InvoiceItemID), job.PatientID, nullUUIDPtr(job.VendorID), job.JobType, string(job.Status),
				nullTimePtr(job.ExpectedDeliveryDate), job.AdvancePayment, nullStringPtr(job.Notes), job.CreatedBy, job.UpdatedBy,
			).
			Suffix("RETURNING id, created_at, updated_at").
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("LabJobRepo.Create build job: %w", err)
		}
		if err := tx.QueryRow(ctx, query, args...).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return fmt.Errorf("LabJobRepo.Create insert job: %w", err)
		}

		return insertLabJobStatusHistory(ctx, tx, clinicID, job.ID, job.Status, nil, job.CreatedBy)
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, clinicID, job.ID)
}

func (r *LabJobRepository) GetByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.LabJobDetails, error) {
	query, args, err := labJobSelect().
		Where(sq.Eq{"lj.id": id, "lj.clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LabJobRepo.GetByID build: %w", err)
