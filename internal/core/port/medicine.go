package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type MedicineRepository interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.Medicine, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Medicine, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Medicine, error)
	Update(ctx context.Context, m *domain.Medicine) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MedicineService interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.Medicine, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Medicine, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Medicine, error)
	Update(ctx context.Context, m *domain.Medicine) error
	Delete(ctx context.Context, id uuid.UUID) error
}
