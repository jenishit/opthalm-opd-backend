package dto

import (
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     *string   `json:"phone"`
	RoleName  string    `json:"role_name"`
	Email     string    `json:"email"`
}

func NewProfileResponse(p *domain.GetProfileDetails) *ProfileResponse {
	return &ProfileResponse{
		ID:        p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		UserID:    p.UserID,
		Phone:     p.Phone,
		RoleName:  p.RoleName,
		Email:     p.Email,
	}
}

func NewProfileResponses(p []*domain.GetProfileDetails) []*ProfileResponse {
	profiles := make([]*ProfileResponse, 0, len(p))

	for _, profile := range p {
		profiles = append(profiles, &ProfileResponse{
			ID:        profile.ID,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			UserID:    profile.UserID,
			Phone:     profile.Phone,
			RoleName:  profile.RoleName,
			Email:     profile.Email,
		})
	}
	return profiles
}

type UpdateProfileRequest struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     *string   `json:"phone"`
	RoleName  string    `json:"role_name"`
	Email     string    `json:"email"`
}
