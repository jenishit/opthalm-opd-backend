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
		return nil, fmt.Errorf("LabJobRepo.listHistory query: %w", err)
	}
	defer rows.Close()

	var history []*domain.LabJobStatusHistory
	for rows.Next() {
		var h domain.LabJobStatusHistory
		var status string
		var notes sql.NullString
		if err := rows.Scan(&h.ID, &h.LabJobID, &status, &h.ChangedAt, &h.ChangedBy, &notes); err != nil {
			return nil, fmt.Errorf("LabJobRepo.listHistory scan: %w", err)
		}
		h.Status = domain.LabJobStatus(status)
		if notes.Valid {
			h.Notes = &notes.String
		}
		history = append(history, &h)
	}
	return history, rows.Err()
}

func labJobSelect() sq.SelectBuilder {
	return sq.Select(
		"lj.id", "lj.invoice_id", "lj.invoice_item_id", "lj.patient_id", "p.full_name", "lj.vendor_id", "v.name",
		"lj.job_type", "lj.status", "lj.expected_delivery_date", "lj.delivered_at", "lj.advance_payment", "lj.notes",
		"lj.created_by", "lj.updated_by", "lj.created_at", "lj.updated_at",
	).
		From("lab_jobs lj").
		Join("patients p ON p.id = lj.patient_id").
		LeftJoin("vendors v ON v.id = lj.vendor_id")
}

func scanLabJobDetailsRow(row pgx.Row) (*domain.LabJobDetails, error) {
	var d domain.LabJobDetails
	var invoiceID, invoiceItemID, vendorID uuid.NullUUID
	var vendorName sql.NullString
	var status string
	var expectedDeliveryDate, deliveredAt sql.NullTime
	var notes sql.NullString

	err := row.Scan(
		&d.ID, &invoiceID, &invoiceItemID, &d.PatientID, &d.PatientName, &vendorID, &vendorName,
		&d.JobType, &status, &expectedDeliveryDate, &deliveredAt, &d.AdvancePayment, &notes,
		&d.CreatedBy, &d.UpdatedBy, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, fmt.Errorf("scan lab job: %w", err)
	}

	d.Status = domain.LabJobStatus(status)
	if invoiceID.Valid {
		d.InvoiceID = &invoiceID.UUID
	}
	if invoiceItemID.Valid {
		d.InvoiceItemID = &invoiceItemID.UUID
	}
	if vendorID.Valid {
		d.VendorID = &vendorID.UUID
	}
	if vendorName.Valid {
