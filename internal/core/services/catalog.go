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
