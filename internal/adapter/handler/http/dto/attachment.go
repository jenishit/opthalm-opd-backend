package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

// UploadAttachmentReq is the non-file form fields that ride alongside the
// "file" part of a multipart/form-data upload — see AttachmentHandler.Upload.
type UploadAttachmentReq struct {
	Category string  `form:"category" binding:"required,oneof=fundus_photo oct visual_field biometry topography other"`
	Eye      string  `form:"eye" binding:"omitempty,oneof=OD OS OU none"`
	Notes    *string `form:"notes"`
}

type AttachmentResponse struct {
	ID               uuid.UUID `json:"id"`
	VisitID          uuid.UUID `json:"visit_id"`
	PatientID        uuid.UUID `json:"patient_id"`
	Category         string    `json:"category"`
	Eye              string    `json:"eye"`
	URL              string    `json:"url"`
	ContentType      string    `json:"content_type"`
	OriginalFilename string    `json:"original_filename,omitempty"`
	SizeBytes        int64     `json:"size_bytes"`
	Notes            *string   `json:"notes,omitempty"`
	UploadedBy       uuid.UUID `json:"uploaded_by"`
	UploadedAt       time.Time `json:"uploaded_at"`
}

// AttachmentRes pairs an attachment's metadata with url — a presigned (or
// reverse-proxied) link to fetch the file, resolved by the caller at
// request time since presigned URLs expire and so are never stored.
func AttachmentRes(a *domain.ClinicalAttachment, url string) *AttachmentResponse {
	if a == nil {
		return nil
	}
	return &AttachmentResponse{
		ID:               a.ID,
		VisitID:          a.VisitID,
		PatientID:        a.PatientID,
		Category:         string(a.Category),
		Eye:              a.Eye,
		URL:              url,
		ContentType:      a.ContentType,
		OriginalFilename: a.OriginalFilename,
		SizeBytes:        a.SizeBytes,
		Notes:            a.Notes,
		UploadedBy:       a.UploadedBy,
		UploadedAt:       a.UploadedAt,
	}
}
