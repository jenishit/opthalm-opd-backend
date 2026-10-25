package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type ReportsService struct {
	repo         port.ReportsRepository
	inventorySvc port.InventoryItemService
}

func NewReportsService(r port.ReportsRepository, inventorySvc port.InventoryItemService) *ReportsService {
	return &ReportsService{repo: r, inventorySvc: inventorySvc}
}

func (s *ReportsService) SalesDaily(ctx context.Context, clinicID uuid.UUID, date time.Time) (*domain.SalesSummary, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	to := from.AddDate(0, 0, 1)
	summary, err := s.repo.SalesSummary(ctx, clinicID, from, to)
	if err != nil {
		return nil, err
	}
	summary.Period = from.Format("2006-01-02")
	return summary, nil
}

func (s *ReportsService) SalesMonthly(ctx context.Context, clinicID uuid.UUID, year, month int) (*domain.SalesSummary, error) {
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)
	summary, err := s.repo.SalesSummary(ctx, clinicID, from, to)
	if err != nil {
		return nil, err
	}
	summary.Period = from.Format("2006-01")
	return summary, nil
}

func (s *ReportsService) SalesRange(ctx context.Context, clinicID uuid.UUID, from, to time.Time) (*domain.SalesSummary, error) {
	summary, err := s.repo.SalesSummary(ctx, clinicID, from, to.AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}
	summary.Period = from.Format("2006-01-02") + " to " + to.Format("2006-01-02")
	return summary, nil
}

func (s *ReportsService) PatientDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.PatientDue, error) {
	return s.repo.PatientDues(ctx, clinicID)
}

func (s *ReportsService) VendorDues(ctx context.Context, clinicID uuid.UUID) ([]*domain.VendorDue, error) {
	return s.repo.VendorDues(ctx, clinicID)
}

func (s *ReportsService) InventoryLowStock(ctx context.Context, clinicID uuid.UUID) ([]*domain.InventoryItem, error) {
	return s.inventorySvc.LowStock(ctx, clinicID)
}

func (s *ReportsService) InventoryValuation(ctx context.Context, clinicID uuid.UUID) (*domain.InventoryValuation, error) {
	return s.repo.InventoryValuation(ctx, clinicID)
}

func (s *ReportsService) VisitsSummaryByDay(ctx context.Context, clinicID uuid.UUID, from, to time.Time) ([]*domain.VisitsSummary, error) {
	return s.repo.VisitsSummaryByDay(ctx, clinicID, from, to.AddDate(0, 0, 1))
}

func (s *ReportsService) VisitsSummaryByDoctor(ctx context.Context, clinicID uuid.UUID, from, to time.Time) ([]*domain.DoctorVisitsSummary, error) {
	return s.repo.VisitsSummaryByDoctor(ctx, clinicID, from, to.AddDate(0, 0, 1))
}
