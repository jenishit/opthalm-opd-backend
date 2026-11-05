-- +goose Up
-- Referral source, captured once at patient registration: who sent this
-- patient to the clinic (a referring doctor/optometrist, or a hospital/
-- clinic name), for referral tracking and reporting.
ALTER TABLE patients ADD COLUMN referred_by_name VARCHAR(255);
ALTER TABLE patients ADD COLUMN referred_by_facility VARCHAR(255);

-- Clinical attachments: images pulled off diagnostic devices (fundus
-- camera, OCT, visual field analyzer, biometry printouts) for a visit.
-- The files themselves live in object storage (MinIO); this table only
-- holds the metadata and the object key.
CREATE TABLE clinical_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    visit_id UUID NOT NULL REFERENCES visits(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    category VARCHAR(30) NOT NULL
        CHECK (category IN ('fundus_photo', 'oct', 'visual_field', 'biometry', 'topography', 'other')),
    eye VARCHAR(4) NOT NULL DEFAULT 'none'
        CHECK (eye IN ('OD', 'OS', 'OU', 'none')),
    object_key VARCHAR(500) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    original_filename VARCHAR(255),
    size_bytes BIGINT NOT NULL,
    notes VARCHAR(500),
    uploaded_by UUID REFERENCES users(id),
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_clinical_attachments_visit_id ON clinical_attachments (visit_id);
CREATE INDEX idx_clinical_attachments_patient_id ON clinical_attachments (patient_id);

-- Consents: a signed consent (exam/procedure/surgery/general) for a visit,
-- with the signature captured as an image (drawn on a signature pad or
-- tablet canvas client-side, exported as PNG/JPEG, and uploaded the same
-- way clinical attachments are) and stored in the same object store.
CREATE TABLE consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    visit_id UUID NOT NULL REFERENCES visits(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    consent_type VARCHAR(30) NOT NULL
        CHECK (consent_type IN ('exam', 'procedure', 'surgery', 'general')),
    signed_by_name VARCHAR(255) NOT NULL,
    signed_by_relation VARCHAR(20) NOT NULL DEFAULT 'self'
        CHECK (signed_by_relation IN ('self', 'guardian', 'other')),
    signature_object_key VARCHAR(500) NOT NULL,
    witnessed_by UUID REFERENCES users(id),
    notes VARCHAR(500),
    created_by UUID REFERENCES users(id),
    signed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_consents_visit_id ON consents (visit_id);
CREATE INDEX idx_consents_patient_id ON consents (patient_id);

-- +goose Down
DROP TABLE consents;
DROP TABLE clinical_attachments;

ALTER TABLE patients DROP COLUMN referred_by_facility;
ALTER TABLE patients DROP COLUMN referred_by_name;
