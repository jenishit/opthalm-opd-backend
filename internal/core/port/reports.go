package port

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ReportsRepository interface {
	SalesSummary(ctx context.Context, clinicID uuid.UUID, from, to time.Time) (*domain.SalesSummary, error)
	PatientDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.PatientDue, error)
	VendorDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.VendorDue, error)
	InventoryValuation(ctx context.Context, clinicID uuid.UUID) (*domain.InventoryValuation, error)
	VisitsSummaryByDay(ctx context.Context, clinicID uuid.UUID, from, to time.Time) ([]*domain.VisitsSummary, error)
	VisitsSummaryByDoctor(ctx context.Context, clinicID uuid.UUID, from, to time.Time) ([]*domain.DoctorVisitsSummary, error)
}

type ReportsService interface {
	SalesDaily(ctx context.Context, clinicID uuid.UUID, date time.Time) (*domain.SalesSummary, error)
	SalesMonthly(ctx context.Context, clinicID uuid.UUID, year, month int) (*domain.SalesSummary, error)
	SalesRange(ctx context.Context, clinicID uuid.UUID, from, to time.Time) (*domain.SalesSummary, error)
	PatientDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.PatientDue, error)
	VendorDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.VendorDue, error)
	InventoryLowStock(ctx context.Context, clinicID uuid.UUID) ([]*domain.InventoryItem, error)
	InventoryValuation(ctx context.Context, clinicID uuid.UUID) (*domain.InventoryValuation, error)
	VisitsSummaryByDay(ctx context.Context, clinicID uuid.UUID, from, to time.Time) ([]*domain.VisitsSummary, error)
	VisitsSummaryByDoctor(ctx context.Context, clinicID uuid.UUID, from, to time.Time) ([]*domain.DoctorVisitsSummary, error)
}
