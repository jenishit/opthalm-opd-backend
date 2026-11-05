package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

// RecordConsentReq is the non-file form fields that ride alongside the
// "signature" part of a multipart/form-data upload — see ConsentHandler.Create.
type RecordConsentReq struct {
	ConsentType      string     `form:"consent_type" binding:"required,oneof=exam procedure surgery general"`
	SignedByName     string     `form:"signed_by_name" binding:"required"`
	SignedByRelation string     `form:"signed_by_relation" binding:"omitempty,oneof=self guardian other"`
	WitnessedBy      *uuid.UUID `form:"witnessed_by"`
	Notes            *string    `form:"notes"`
}

type ConsentResponse struct {
	ID               uuid.UUID  `json:"id"`
	VisitID          uuid.UUID  `json:"visit_id"`
	PatientID        uuid.UUID  `json:"patient_id"`
	ConsentType      string     `json:"consent_type"`
	SignedByName     string     `json:"signed_by_name"`
	SignedByRelation string     `json:"signed_by_relation"`
	SignatureURL     string     `json:"signature_url"`
	WitnessedBy      *uuid.UUID `json:"witnessed_by,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
	CreatedBy        uuid.UUID  `json:"created_by"`
	SignedAt         time.Time  `json:"signed_at"`
}

func ConsentRes(c *domain.Consent, signatureURL string) *ConsentResponse {
	if c == nil {
		return nil
	}
	return &ConsentResponse{
		ID:               c.ID,
		VisitID:          c.VisitID,
		PatientID:        c.PatientID,
		ConsentType:      string(c.ConsentType),
		SignedByName:     c.SignedByName,
		SignedByRelation: string(c.SignedByRelation),
		SignatureURL:     signatureURL,
		WitnessedBy:      c.WitnessedBy,
		Notes:            c.Notes,
		CreatedBy:        c.CreatedBy,
		SignedAt:         c.SignedAt,
	}
}
