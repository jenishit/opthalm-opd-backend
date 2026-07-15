package dto

import (
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type CreateRole struct { //binding required is the compulsory values
	RoleName string `json:"role_name" binding:"required"`
}

type Role struct {
	RoleID   uuid.UUID `json:"role_id"`
	RoleName string    `json:"role_name"`
}

func RoleResponse(r *domain.Role) *Role {
	return &Role{
		RoleID:   r.ID,
		RoleName: r.RoleName,
	}
}

// func RolesResponse(r *domain.Role) []*Role {
// 	roles := 
// }