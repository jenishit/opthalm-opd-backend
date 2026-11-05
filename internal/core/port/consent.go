package port

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type ConsentRepository interface {
	CreateConsent(ctx context.Context, clinicID uuid.UUID, c *domain.Consent) (*domain.Consent, error)
	ListConsentsByVisitID(ctx context.Context, clinicID, visitID uuid.UUID) ([]*domain.Consent, error)
}

type ConsentService interface {
	RecordConsent(
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
	) (*domain.Consent, string, error)
	ListConsents(ctx context.Context, clinicID, visitID uuid.UUID) ([]*ConsentWithURL, error)
}

// ConsentWithURL pairs a consent record with a short-lived URL to view the
// signature image.
type ConsentWithURL struct {
	*domain.Consent
	SignatureURL string
}
