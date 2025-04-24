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
