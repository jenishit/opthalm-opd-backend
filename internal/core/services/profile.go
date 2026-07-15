package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type ProfileService struct {
	repo port.ProfileRepository
}

func NewProfileService(pr port.ProfileRepository) *ProfileService {
	return &ProfileService{
		repo: pr,
	}
}

func (p *ProfileService) CreateProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	return p.repo.CreateProfile(ctx, profile)
}

func (p *ProfileService) GetProfileByID(ctx context.Context, id uuid.UUID) (*domain.GetProfileDetails, error) {
	profile, err := p.repo.GetProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return profile, err
}

func (p *ProfileService) GetProfiles(ctx context.Context) ([]*domain.GetProfileDetails, error) {
	profiles, err := p.repo.GetProfiles(ctx)
	if err != nil {
		return nil, err
	}

	return profiles, err
}

func (p *ProfileService) UpdateProfileByUserID(ctx context.Context, prof *domain.GetProfileDetails) error {
	if prof.UserID == uuid.Nil {
		return errors.New("user id is required")
	}
	return p.repo.UpdateProfileByUserID(ctx, prof)
}
