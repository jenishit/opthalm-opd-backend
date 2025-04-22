package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type MedicineService struct {
	repo port.MedicineRepository
}

func NewMedicineService(r port.MedicineRepository) *MedicineService {
	return &MedicineService{repo: r}
}

func (s *MedicineService) Search(ctx context.Context, query string, limit int) ([]*domain.Medicine, error) {
	return s.repo.Search(ctx, query, limit)
}

func (s *MedicineService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Medicine, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MedicineService) List(ctx context.Context, limit, offset int) ([]*domain.Medicine, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *MedicineService) Update(ctx context.Context, m *domain.Medicine) error {
	return s.repo.Update(ctx, m)
}

func (s *MedicineService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

type DiagnosisCatalogService struct {
	repo port.DiagnosisCatalogRepository
}

func NewDiagnosisCatalogService(r port.DiagnosisCatalogRepository) *DiagnosisCatalogService {
	return &DiagnosisCatalogService{repo: r}
}

func (s *DiagnosisCatalogService) Search(ctx context.Context, query string, limit int) ([]*domain.DiagnosisCatalog, error) {
	return s.repo.Search(ctx, query, limit)
}

func (s *DiagnosisCatalogService) GetByID(ctx context.Context, id uuid.UUID) (*domain.DiagnosisCatalog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DiagnosisCatalogService) List(ctx context.Context, limit, offset int) ([]*domain.DiagnosisCatalog, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *DiagnosisCatalogService) Update(ctx context.Context, d *domain.DiagnosisCatalog) error {
	return s.repo.Update(ctx, d)
}

func (s *DiagnosisCatalogService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

type HistoryConditionService struct {
	repo port.HistoryConditionRepository
}

func NewHistoryConditionService(r port.HistoryConditionRepository) *HistoryConditionService {
	return &HistoryConditionService{repo: r}
}

func (s *HistoryConditionService) Search(ctx context.Context, query string, limit int) ([]*domain.HistoryCondition, error) {
	return s.repo.Search(ctx, query, limit)
}

func (s *HistoryConditionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.HistoryCondition, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HistoryConditionService) List(ctx context.Context, limit, offset int) ([]*domain.HistoryCondition, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *HistoryConditionService) Update(ctx context.Context, h *domain.HistoryCondition) error {
	return s.repo.Update(ctx, h)
}

func (s *HistoryConditionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
