package services

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type AttachmentService struct {
	repo    port.AttachmentRepository
	objects port.ObjectStorage
}

func NewAttachmentService(repo port.AttachmentRepository, objects port.ObjectStorage) *AttachmentService {
	return &AttachmentService{repo: repo, objects: objects}
}

// UploadAttachment stores the file bytes in object storage under a key
// namespaced by clinic/visit (so one clinic's files are never reachable
// through another's key, and a visit's files are grouped together), then
// records the metadata row. If the metadata write fails, the just-uploaded
// object is cleaned up rather than left orphaned.
func (s *AttachmentService) UploadAttachment(
	ctx context.Context,
	clinicID, visitID, patientID, uploadedBy uuid.UUID,
	category domain.AttachmentCategory,
	eye string,
	filename, contentType string,
	size int64,
	r io.Reader,
	notes *string,
) (*domain.ClinicalAttachment, string, error) {
	objectKey := fmt.Sprintf("clinics/%s/visits/%s/%s-%s", clinicID, visitID, uuid.NewString(), filename)

	if err := s.objects.Upload(ctx, objectKey, r, size, contentType); err != nil {
		return nil, "", fmt.Errorf("uploading attachment: %w", err)
	}

	attachment := &domain.ClinicalAttachment{
		VisitID:          visitID,
		PatientID:        patientID,
		Category:         category,
		Eye:              eye,
		ObjectKey:        objectKey,
		ContentType:      contentType,
		OriginalFilename: filename,
		SizeBytes:        size,
		Notes:            notes,
		UploadedBy:       uploadedBy,
	}

	saved, err := s.repo.CreateAttachment(ctx, clinicID, attachment)
	if err != nil {
		_ = s.objects.Delete(ctx, objectKey)
		return nil, "", fmt.Errorf("saving attachment metadata: %w", err)
	}

	url, err := s.objects.PresignedURL(ctx, objectKey, port.AttachmentURLExpiry)
	if err != nil {
		return nil, "", fmt.Errorf("presigning uploaded attachment: %w", err)
	}

	return saved, url, nil
}

func (s *AttachmentService) ListAttachments(ctx context.Context, clinicID, visitID uuid.UUID) ([]*port.AttachmentWithURL, error) {
	attachments, err := s.repo.ListAttachmentsByVisitID(ctx, clinicID, visitID)
	if err != nil {
		return nil, err
	}

	result := make([]*port.AttachmentWithURL, 0, len(attachments))
	for _, a := range attachments {
		url, err := s.objects.PresignedURL(ctx, a.ObjectKey, port.AttachmentURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("presigning attachment %s: %w", a.ID, err)
		}
		result = append(result, &port.AttachmentWithURL{ClinicalAttachment: a, URL: url})
	}

	return result, nil
}

func (s *AttachmentService) DeleteAttachment(ctx context.Context, clinicID, id uuid.UUID) error {
	attachment, err := s.repo.GetAttachmentByID(ctx, clinicID, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteAttachmentByID(ctx, clinicID, id); err != nil {
		return err
	}

	// Best-effort: the metadata row is already gone either way, so a
	// failure here just leaves an unreferenced object behind rather than
	// blocking the delete the caller asked for.
	_ = s.objects.Delete(ctx, attachment.ObjectKey)
	return nil
}
