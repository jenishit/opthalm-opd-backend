package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConsentType string

const (
	ConsentExam      ConsentType = "exam"
	ConsentProcedure ConsentType = "procedure"
	ConsentSurgery   ConsentType = "surgery"
	ConsentGeneral   ConsentType = "general"
)

type SignedByRelation string

const (
	SignedBySelf     SignedByRelation = "self"
	SignedByGuardian SignedByRelation = "guardian"
	SignedByOther    SignedByRelation = "other"
)

// Consent records a signed consent for a visit. The signature itself is
// captured as an image — drawn on a signature pad/tablet canvas client-side
// and exported as PNG/JPEG — uploaded and stored in object storage the same
// way a ClinicalAttachment is; SignatureObjectKey points to it there.
type Consent struct {
	ID                 uuid.UUID
	ClinicID           uuid.UUID
	VisitID            uuid.UUID
	PatientID          uuid.UUID
	ConsentType        ConsentType
	SignedByName       string
	SignedByRelation   SignedByRelation
	SignatureObjectKey string
	WitnessedBy        *uuid.UUID
	Notes              *string
	CreatedBy          uuid.UUID
	SignedAt           time.Time
}
