package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type DiagnosisCatalogRepository interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.DiagnosisCatalog, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.DiagnosisCatalog, error)
	List(ctx context.Context, limit, offset int) ([]*domain.DiagnosisCatalog, error)
	Update(ctx context.Context, d *domain.DiagnosisCatalog) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DiagnosisCatalogService interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.DiagnosisCatalog, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.DiagnosisCatalog, error)
	List(ctx context.Context, limit, offset int) ([]*domain.DiagnosisCatalog, error)
	Update(ctx context.Context, d *domain.DiagnosisCatalog) error
	Delete(ctx context.Context, id uuid.UUID) error
}
