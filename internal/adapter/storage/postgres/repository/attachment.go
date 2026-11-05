package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type AttachmentRepository struct {
	DB *postgres.DB
}

func NewAttachmentRepository(db *postgres.DB) *AttachmentRepository {
	return &AttachmentRepository{DB: db}
}

func (r *AttachmentRepository) CreateAttachment(ctx context.Context, clinicID uuid.UUID, a *domain.ClinicalAttachment) (*domain.ClinicalAttachment, error) {
	query, args, err := sq.
		Insert("clinical_attachments").
		Columns(
			"clinic_id",
			"visit_id",
			"patient_id",
			"category",
			"eye",
			"object_key",
			"content_type",
			"original_filename",
			"size_bytes",
			"notes",
			"uploaded_by",
		).
		Values(
			clinicID,
			a.VisitID,
			a.PatientID,
			a.Category,
			a.Eye,
			a.ObjectKey,
			a.ContentType,
			nullString(a.OriginalFilename),
			a.SizeBytes,
			a.Notes,
			a.UploadedBy,
		).
		Suffix("RETURNING id, uploaded_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AttachmentRepo.CreateAttachment build: %w", err)
	}

	if err := r.DB.QueryRow(ctx, query, args...).Scan(&a.ID, &a.UploadedAt); err != nil {
		return nil, fmt.Errorf("AttachmentRepo.CreateAttachment scan: %w", err)
	}

	a.ClinicID = clinicID
	return a, nil
}

func (r *AttachmentRepository) ListAttachmentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.ClinicalAttachment, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"patient_id",
			"category",
			"eye",
			"object_key",
			"content_type",
			"original_filename",
			"size_bytes",
			"notes",
			"uploaded_by",
			"uploaded_at",
		).
		From("clinical_attachments").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		OrderBy("uploaded_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("AttachmentRepo.ListByVisitID query: %w", err)
	}
	defer rows.Close()

	attachments := make([]*domain.ClinicalAttachment, 0)
	for rows.Next() {
		a, err := scanAttachment(rows)
		if err != nil {
			return nil, err
		}
		a.ClinicID = clinicID
		attachments = append(attachments, a)
	}

	return attachments, nil
}

func (r *AttachmentRepository) GetAttachmentByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.ClinicalAttachment, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"patient_id",
			"category",
			"eye",
			"object_key",
			"content_type",
			"original_filename",
			"size_bytes",
			"notes",
			"uploaded_by",
			"uploaded_at",
		).
		From("clinical_attachments").
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	a, err := scanAttachment(r.DB.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}
	a.ClinicID = clinicID
	return a, nil
}

func (r *AttachmentRepository) DeleteAttachmentByID(ctx context.Context, clinicID, id uuid.UUID) error {
	query, args, err := sq.
		Delete("clinical_attachments").
		Where(sq.Eq{"id": id, "clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("AttachmentRepo.DeleteByID build: %w", err)
	}

	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("AttachmentRepo.DeleteByID exec: %w", err)
	}
	return nil
}

// rowScanner is satisfied by both pgx.Row and pgx.Rows, so scanAttachment
// works for both the single-row and multi-row read paths above.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanAttachment(row rowScanner) (*domain.ClinicalAttachment, error) {
	var a domain.ClinicalAttachment
	var originalFilename, notes sql.NullString

	err := row.Scan(
		&a.ID,
		&a.VisitID,
		&a.PatientID,
		&a.Category,
		&a.Eye,
		&a.ObjectKey,
		&a.ContentType,
		&originalFilename,
		&a.SizeBytes,
		&notes,
		&a.UploadedBy,
		&a.UploadedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AttachmentRepo scan: %w", err)
	}

	if originalFilename.Valid {
		a.OriginalFilename = originalFilename.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}

	return &a, nil
}
