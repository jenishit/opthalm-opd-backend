package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
	GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.GetProfileDetails, error)
	GetProfiles(ctx context.Context) ([]*domain.GetProfileDetails, error)
	UpdateProfileByUserID(ctx context.Context, prof *domain.GetProfileDetails) error 
}

type ProfileService interface {
	CreateProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
	GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.GetProfileDetails, error)
	GetProfiles(ctx context.Context) ([]*domain.GetProfileDetails, error)
	UpdateProfileByUserID(ctx context.Context, prof *domain.GetProfileDetails) error 
}