package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type CreatePatientHistoryReq struct {
	UsesSpectacles      bool        `json:"uses_spectacles"`
	SystemicMedications *string     `json:"systemic_medications"`
	TopicalMedications  *string     `json:"topical_medications"`
	Notes               *string     `json:"notes"`
	ConditionIDs        []uuid.UUID `json:"condition_ids" binding:"omitempty,dive"`
}

type PatientHistoryResponse struct {
	ID                  uuid.UUID   `json:"id"`
	PatientID           uuid.UUID   `json:"patient_id"`
	UsesSpectacles      bool        `json:"uses_spectacles"`
	SystemicMedications *string     `json:"systemic_medications,omitempty"`
	TopicalMedications  *string     `json:"topical_medications,omitempty"`
	Notes               *string     `json:"notes,omitempty"`
	ConditionIDs        []uuid.UUID `json:"condition_ids"`
	CreatedAt           time.Time   `json:"created_at"`
	CreatedBy           uuid.UUID   `json:"created_by,omitempty"`
}

func PatientHistoryRes(h *domain.PatientHistory) *PatientHistoryResponse {
	if h == nil {
		return nil
	}

	conditionIDs := h.ConditionIDs
	if conditionIDs == nil {
		conditionIDs = []uuid.UUID{}
	}

	return &PatientHistoryResponse{
		ID:                  h.ID,
		PatientID:           h.PatientID,
		UsesSpectacles:      h.UsesSpectacles,
		SystemicMedications: h.SystemicMedications,
		TopicalMedications:  h.TopicalMedications,
		Notes:               h.Notes,
		ConditionIDs:        conditionIDs,
		CreatedAt:           h.CreatedAt,
		CreatedBy:           h.CreatedBy,
	}
}

func PatientHistoryResList(history []*domain.PatientHistory) []*PatientHistoryResponse {
	rsp := make([]*PatientHistoryResponse, 0, len(history))
	for _, h := range history {
		rsp = append(rsp, PatientHistoryRes(h))
	}
	return rsp
}
