package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type CreateClinicLocationReq struct {
	Name      string  `json:"name" binding:"required"`
	Address   *string `json:"address"`
	ExamRooms int     `json:"exam_rooms" binding:"min=0"`
	Status    string  `json:"status" binding:"omitempty,oneof=active opening_soon"`
}

type ClinicLocationResponse struct {
	ID        uuid.UUID `json:"id"`
	ClinicID  uuid.UUID `json:"clinic_id"`
	Name      string    `json:"name"`
	Address   *string   `json:"address,omitempty"`
	ExamRooms int       `json:"exam_rooms"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func ClinicLocationRes(l *domain.ClinicLocation) *ClinicLocationResponse {
	if l == nil {
		return nil
	}

	return &ClinicLocationResponse{
		ID:        l.ID,
		ClinicID:  l.ClinicID,
		Name:      l.Name,
		Address:   l.Address,
		ExamRooms: l.ExamRooms,
		Status:    l.Status,
		CreatedAt: l.CreatedAt,
	}
}

func ClinicLocationResList(locations []*domain.ClinicLocation) []*ClinicLocationResponse {
	rsp := make([]*ClinicLocationResponse, 0, len(locations))
	for _, l := range locations {
		rsp = append(rsp, ClinicLocationRes(l))
	}
	return rsp
}
