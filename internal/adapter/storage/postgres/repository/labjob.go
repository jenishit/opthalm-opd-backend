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
	}

	details, err := scanLabJobDetailsRow(r.DB.QueryRow(ctx, query, args...))
	if err != nil {
		return nil, err
	}

	history, err := r.listHistory(ctx, id)
	if err != nil {
		return nil, err
	}
	details.History = history

	return details, nil
}

func (r *LabJobRepository) List(ctx context.Context, clinicID uuid.UUID, limit, offset int) ([]*domain.LabJobDetails, error) {
	query, args, err := labJobSelect().
		Where(sq.Eq{"lj.clinic_id": clinicID}).
		OrderBy("lj.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LabJobRepo.List build: %w", err)
	}
	return scanLabJobDetailsList(ctx, r.DB, query, args)
}

func (r *LabJobRepository) ListByPatientID(ctx context.Context, clinicID, patientID uuid.UUID) ([]*domain.LabJobDetails, error) {
	query, args, err := labJobSelect().
		Where(sq.Eq{"lj.patient_id": patientID, "lj.clinic_id": clinicID}).
		OrderBy("lj.created_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LabJobRepo.ListByPatientID build: %w", err)
	}
	return scanLabJobDetailsList(ctx, r.DB, query, args)
}

func (r *LabJobRepository) UpdateStatus(ctx context.Context, clinicID, id uuid.UUID, status domain.LabJobStatus, notes *string, changedBy uuid.UUID) error {
	return r.DB.WithTx(ctx, func(tx pgx.Tx) error {
		ub := sq.Update("lab_jobs").
			Set("status", string(status)).
			Set("updated_by", changedBy).
			Set("updated_at", sq.Expr("NOW()"))
		if status == domain.LabJobDelivered {
			ub = ub.Set("delivered_at", sq.Expr("NOW()"))
		}
		query, args, err := ub.Where(sq.Eq{"id": id, "clinic_id": clinicID}).PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("LabJobRepo.UpdateStatus build: %w", err)
		}
		tag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("LabJobRepo.UpdateStatus exec: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrDataNotFound
		}

		return insertLabJobStatusHistory(ctx, tx, clinicID, id, status, notes, changedBy)
	})
}

func insertLabJobStatusHistory(ctx context.Context, tx pgx.Tx, clinicID, jobID uuid.UUID, status domain.LabJobStatus, notes *string, changedBy uuid.UUID) error {
	query, args, err := sq.Insert("lab_job_status_history").
		Columns("clinic_id", "lab_job_id", "status", "changed_by", "notes").
		Values(clinicID, jobID, string(status), changedBy, nullStringPtr(notes)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("insert lab job status history build: %w", err)
	}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert lab job status history exec: %w", err)
	}
	return nil
}

func (r *LabJobRepository) listHistory(ctx context.Context, jobID uuid.UUID) ([]*domain.LabJobStatusHistory, error) {
	query, args, err := sq.Select("id", "lab_job_id", "status", "changed_at", "changed_by", "notes").
		From("lab_job_status_history").
		Where(sq.Eq{"lab_job_id": jobID}).
		OrderBy("changed_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("LabJobRepo.listHistory build: %w", err)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
