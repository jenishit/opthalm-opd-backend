package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ConsentRepository struct {
	DB *postgres.DB
}

func NewConsentRepository(db *postgres.DB) *ConsentRepository {
	return &ConsentRepository{DB: db}
}

func (r *ConsentRepository) CreateConsent(ctx context.Context, clinicID uuid.UUID, c *domain.Consent) (*domain.Consent, error) {
	query, args, err := sq.
		Insert("consents").
		Columns(
			"clinic_id",
			"visit_id",
			"patient_id",
			"consent_type",
			"signed_by_name",
			"signed_by_relation",
			"signature_object_key",
			"witnessed_by",
			"notes",
			"created_by",
		).
		Values(
			clinicID,
			c.VisitID,
			c.PatientID,
			c.ConsentType,
			c.SignedByName,
			c.SignedByRelation,
			c.SignatureObjectKey,
			nullUUIDPtr(c.WitnessedBy),
			c.Notes,
			c.CreatedBy,
		).
		Suffix("RETURNING id, signed_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ConsentRepo.CreateConsent build: %w", err)
	}

	if err := r.DB.QueryRow(ctx, query, args...).Scan(&c.ID, &c.SignedAt); err != nil {
		return nil, fmt.Errorf("ConsentRepo.CreateConsent scan: %w", err)
	}

	c.ClinicID = clinicID
	return c, nil
}

func (r *ConsentRepository) ListConsentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.Consent, error) {
	query, args, err := sq.
		Select(
			"id",
			"visit_id",
			"patient_id",
			"consent_type",
			"signed_by_name",
			"signed_by_relation",
			"signature_object_key",
			"witnessed_by",
			"notes",
			"created_by",
			"signed_at",
		).
		From("consents").
		Where(sq.Eq{"visit_id": visitID, "clinic_id": clinicID}).
		OrderBy("signed_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ConsentRepo.ListByVisitID query: %w", err)
	}
	defer rows.Close()

	consents := make([]*domain.Consent, 0)
	for rows.Next() {
		var c domain.Consent
		var witnessedBy uuid.NullUUID
		var notes sql.NullString

		if err := rows.Scan(
			&c.ID,
			&c.VisitID,
			&c.PatientID,
			&c.ConsentType,
			&c.SignedByName,
			&c.SignedByRelation,
			&c.SignatureObjectKey,
			&witnessedBy,
			&notes,
			&c.CreatedBy,
			&c.SignedAt,
		); err != nil {
			return nil, fmt.Errorf("ConsentRepo.ListByVisitID scan: %w", err)
		}

		if witnessedBy.Valid {
			c.WitnessedBy = &witnessedBy.UUID
		}
		if notes.Valid {
			c.Notes = &notes.String
		}

		c.ClinicID = clinicID
		consents = append(consents, &c)
	}

	return consents, nil
}
