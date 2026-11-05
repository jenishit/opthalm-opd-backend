package port

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type AttachmentRepository interface {
	CreateAttachment(ctx context.Context, clinicID uuid.UUID, a *domain.ClinicalAttachment) (*domain.ClinicalAttachment, error)
	ListAttachmentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.ClinicalAttachment, error)
	GetAttachmentByID(ctx context.Context, clinicID, id uuid.UUID) (*domain.ClinicalAttachment, error)
	DeleteAttachmentByID(ctx context.Context, clinicID, id uuid.UUID) error
}

// AttachmentService orchestrates uploading the file bytes to object storage
// and recording the metadata row, so a failed upload never leaves a
// dangling DB row and a failed DB write never leaves an orphaned object.
type AttachmentService interface {
	UploadAttachment(ctx context.Context, clinicID, visitID, patientID, uploadedBy uuid.UUID, category domain.AttachmentCategory, eye string, filename, contentType string, size int64, r io.Reader, notes *string) (*domain.ClinicalAttachment, string, error)
	ListAttachments(ctx context.Context, clinicID, visitID uuid.UUID) ([]*AttachmentWithURL, error)
	DeleteAttachment(ctx context.Context, clinicID, id uuid.UUID) error
}

// AttachmentWithURL pairs an attachment's metadata with a short-lived URL
// to fetch the actual file — resolved at list time rather than stored,
// since presigned URLs expire.
type AttachmentWithURL struct {
	*domain.ClinicalAttachment
	URL string
}

const AttachmentURLExpiry = 15 * time.Minute
