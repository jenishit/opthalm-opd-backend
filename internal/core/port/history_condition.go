package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type HistoryConditionRepository interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.HistoryCondition, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.HistoryCondition, error)
	List(ctx context.Context, limit, offset int) ([]*domain.HistoryCondition, error)
	Update(ctx context.Context, h *domain.HistoryCondition) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type HistoryConditionService interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.HistoryCondition, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.HistoryCondition, error)
	List(ctx context.Context, limit, offset int) ([]*domain.HistoryCondition, error)
	Update(ctx context.Context, h *domain.HistoryCondition) error
	Delete(ctx context.Context, id uuid.UUID) error
}
