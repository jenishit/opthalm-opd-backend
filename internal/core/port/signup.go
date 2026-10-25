package port

import (
	"context"

	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

// SignupRepository creates a brand-new tenant (clinic + its first admin user
// + a trialing subscription) atomically. It exists to break the chicken-and-
// egg problem of clinic creation being admin-gated while a clinic's first
// admin can't exist without a clinic.
type SignupRepository interface {
	CreateTenant(ctx context.Context, req *domain.SignupRequest, passwordHash string) (*domain.BasicDetails, error)
}
