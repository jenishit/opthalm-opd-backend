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
