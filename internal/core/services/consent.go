package services

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

type ConsentService struct {
	repo    port.ConsentRepository
	objects port.ObjectStorage
}

func NewConsentService(repo port.ConsentRepository, objects port.ObjectStorage) *ConsentService {
	return &ConsentService{repo: repo, objects: objects}
}

// RecordConsent uploads the signature image (drawn on a signature pad or
// tablet canvas client-side and exported as PNG/JPEG) to object storage
// and records the consent. Mirrors AttachmentService.UploadAttachment's
// cleanup-on-failure behavior.
func (s *ConsentService) RecordConsent(
	ctx context.Context,
	clinicID, visitID, patientID, createdBy uuid.UUID,
	consentType domain.ConsentType,
	signedByName string,
	signedByRelation domain.SignedByRelation,
	witnessedBy *uuid.UUID,
	notes *string,
	signatureContentType string,
	signatureSize int64,
	signature io.Reader,
) (*domain.Consent, string, error) {
	objectKey := fmt.Sprintf("clinics/%s/visits/%s/consent-%s.png", clinicID, visitID, uuid.NewString())

	if err := s.objects.Upload(ctx, objectKey, signature, signatureSize, signatureContentType); err != nil {
		return nil, "", fmt.Errorf("uploading signature: %w", err)
	}

	consent := &domain.Consent{
		VisitID:            visitID,
		PatientID:          patientID,
		ConsentType:        consentType,
		SignedByName:       signedByName,
		SignedByRelation:   signedByRelation,
		SignatureObjectKey: objectKey,
		WitnessedBy:        witnessedBy,
		Notes:              notes,
		CreatedBy:          createdBy,
	}

	saved, err := s.repo.CreateConsent(ctx, clinicID, consent)
	if err != nil {
		_ = s.objects.Delete(ctx, objectKey)
		return nil, "", fmt.Errorf("saving consent: %w", err)
	}

	url, err := s.objects.PresignedURL(ctx, objectKey, port.AttachmentURLExpiry)
	if err != nil {
		return nil, "", fmt.Errorf("presigning signature: %w", err)
	}

	return saved, url, nil
}

func (s *ConsentService) ListConsents(ctx context.Context, clinicID, visitID uuid.UUID) ([]*port.ConsentWithURL, error) {
	consents, err := s.repo.ListConsentsByVisitID(ctx, clinicID, visitID)
	if err != nil {
		return nil, err
	}

	result := make([]*port.ConsentWithURL, 0, len(consents))
	for _, c := range consents {
		url, err := s.objects.PresignedURL(ctx, c.SignatureObjectKey, port.AttachmentURLExpiry)
		if err != nil {
			return nil, fmt.Errorf("presigning consent signature %s: %w", c.ID, err)
		}
		result = append(result, &port.ConsentWithURL{Consent: c, SignatureURL: url})
	}

	return result, nil
}
