package domain

import (
	"time"

	"github.com/google/uuid"
)

type AttachmentCategory string

const (
	AttachmentFundusPhoto  AttachmentCategory = "fundus_photo"
	AttachmentOCT          AttachmentCategory = "oct"
	AttachmentVisualField  AttachmentCategory = "visual_field"
	AttachmentBiometry     AttachmentCategory = "biometry"
	AttachmentTopography   AttachmentCategory = "topography"
	AttachmentOtherImaging AttachmentCategory = "other"
)

// ClinicalAttachment is one image file for a visit — a fundus photo, OCT
// scan export, visual field printout, etc. pulled off a diagnostic device
// (downloaded from the device onto a computer/tablet, then uploaded here).
// The bytes live in object storage (MinIO); ObjectKey is the pointer to
// them there.
type ClinicalAttachment struct {
	ID               uuid.UUID
	ClinicID         uuid.UUID
	VisitID          uuid.UUID
	PatientID        uuid.UUID
	Category         AttachmentCategory
	Eye              string // OD | OS | OU | none
	ObjectKey        string
	ContentType      string
	OriginalFilename string
	SizeBytes        int64
	Notes            *string
	UploadedBy       uuid.UUID
	UploadedAt       time.Time
}
