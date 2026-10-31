package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ClinicLocationRepository struct {
	DB *postgres.DB
}

func NewClinicLocationRepository(db *postgres.DB) *ClinicLocationRepository {
	return &ClinicLocationRepository{DB: db}
}

func (r *ClinicLocationRepository) CreateLocation(ctx context.Context, clinicID uuid.UUID, l *domain.ClinicLocation) (*domain.ClinicLocation, error) {
	status := l.Status
	if status == "" {
		status = "active"
	}

	query, args, err := sq.
		Insert("clinic_locations").
		Columns("clinic_id", "name", "address", "exam_rooms", "status", "created_by").
		Values(clinicID, l.Name, l.Address, l.ExamRooms, status, l.CreatedBy).
		Suffix("RETURNING id, clinic_id, status, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ClinicLocation.Create build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(&l.ID, &l.ClinicID, &l.Status, &l.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("ClinicLocation.Create scan: %w", err)
	}

	return l, nil
}

func (r *ClinicLocationRepository) ListLocations(ctx context.Context, clinicID uuid.UUID) ([]*domain.ClinicLocation, error) {
	query, args, err := sq.
		Select("id", "clinic_id", "name", "address", "exam_rooms", "status", "created_by", "created_at").
		From("clinic_locations").
		Where(sq.Eq{"clinic_id": clinicID}).
		OrderBy("created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ClinicLocation.List query: %w", err)
	}
	defer rows.Close()

	locations := make([]*domain.ClinicLocation, 0)

	for rows.Next() {
		var l domain.ClinicLocation

		if err := rows.Scan(&l.ID, &l.ClinicID, &l.Name, &l.Address, &l.ExamRooms, &l.Status, &l.CreatedBy, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("ClinicLocation.List scan: %w", err)
		}

		locations = append(locations, &l)
	}

	return locations, nil
}
