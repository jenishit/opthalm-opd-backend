-- +goose Up
-- visit_symptoms was missing the unique constraint its own upsert
-- (ON CONFLICT (visit_id)) depends on — one symptom row per visit.
ALTER TABLE visit_symptoms ADD CONSTRAINT visit_symptoms_visit_id_key UNIQUE (visit_id);

-- headache and double_vision were modelled as plain booleans, but the
-- symptom chart treats every symptom uniformly as a laterality enum
-- ("none"/"OD"/"OS"/"OU"), so bring these two in line with the rest.
ALTER TABLE visit_symptoms ALTER COLUMN headache DROP DEFAULT;
ALTER TABLE visit_symptoms ALTER COLUMN headache TYPE VARCHAR(4) USING (CASE WHEN headache THEN 'OU' ELSE 'none' END);
ALTER TABLE visit_symptoms ALTER COLUMN headache SET DEFAULT 'none';
ALTER TABLE visit_symptoms ALTER COLUMN headache SET NOT NULL;

ALTER TABLE visit_symptoms ALTER COLUMN double_vision DROP DEFAULT;
ALTER TABLE visit_symptoms ALTER COLUMN double_vision TYPE VARCHAR(4) USING (CASE WHEN double_vision THEN 'OU' ELSE 'none' END);
ALTER TABLE visit_symptoms ALTER COLUMN double_vision SET DEFAULT 'none';
ALTER TABLE visit_symptoms ALTER COLUMN double_vision SET NOT NULL;

ALTER TABLE examination_findings ADD COLUMN id UUID DEFAULT gen_random_uuid();
UPDATE examination_findings SET id = gen_random_uuid() WHERE id IS NULL;
ALTER TABLE examination_findings ALTER COLUMN id SET NOT NULL;
ALTER TABLE examination_findings ADD PRIMARY KEY (id);
ALTER TABLE examination_findings ADD COLUMN iop_time VARCHAR(20);
ALTER TABLE examination_findings ADD COLUMN bp VARCHAR(20);
ALTER TABLE examination_findings ADD COLUMN updated_by UUID REFERENCES users(id);
ALTER TABLE examination_findings ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE examination_findings ADD CONSTRAINT examination_findings_visit_id_key UNIQUE (visit_id);

ALTER TABLE refraction_readings ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE refraction_readings ADD CONSTRAINT refraction_readings_visit_eye_stage_key UNIQUE (visit_id, eye, stage);
ALTER TABLE refraction_readings ADD CONSTRAINT refraction_readings_eye_check
    CHECK (eye IN ('OD', 'OS'));
ALTER TABLE refraction_readings ADD CONSTRAINT refraction_readings_stage_check
    CHECK (stage IN ('unaided', 'old_rx', 'pinhole', 'retinoscopy', 'subjective_acceptance', 'final_rx'));
ALTER TABLE refraction_readings ADD CONSTRAINT refraction_readings_lens_type_check
    CHECK (lens_type IS NULL OR lens_type IN ('single_vision', 'bifocal', 'progressive', 'contact_lens', 'none'));

ALTER TABLE follow_ups ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE follow_ups ADD CONSTRAINT follow_ups_visit_id_key UNIQUE (visit_id);

CREATE TABLE visit_diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    visit_id UUID NOT NULL REFERENCES visits(id),
    diagnosis_id UUID NOT NULL,
    diagnosis_name VARCHAR(255),
    icd10_code VARCHAR(10),
    laterality VARCHAR(4) NOT NULL DEFAULT 'none'
        CHECK (laterality IN ('OD', 'OS', 'OU', 'none')),
    status VARCHAR(20)
        CHECK (status IS NULL OR status IN ('provisional', 'confirmed', 'resolved')),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_visit_diagnoses_visit_id ON visit_diagnoses (visit_id);

CREATE TABLE treatments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    visit_id UUID NOT NULL REFERENCES visits(id),
    type VARCHAR(20) NOT NULL
        CHECK (type IN ('medicine', 'procedure', 'advice')),
    medicine_id UUID REFERENCES medicines(id),
    medicine_name VARCHAR(255),
    procedure_name VARCHAR(255),
    advice_text VARCHAR(500),
    route VARCHAR(50),
    dosage VARCHAR(50),
    frequency VARCHAR(50),
    duration VARCHAR(50),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_treatments_visit_id ON treatments (visit_id);

-- +goose Down
DROP TABLE treatments;
DROP TABLE visit_diagnoses;

ALTER TABLE follow_ups DROP CONSTRAINT follow_ups_visit_id_key;
ALTER TABLE follow_ups DROP COLUMN created_at;

ALTER TABLE refraction_readings DROP CONSTRAINT refraction_readings_lens_type_check;
ALTER TABLE refraction_readings DROP CONSTRAINT refraction_readings_stage_check;
ALTER TABLE refraction_readings DROP CONSTRAINT refraction_readings_eye_check;
ALTER TABLE refraction_readings DROP CONSTRAINT refraction_readings_visit_eye_stage_key;
ALTER TABLE refraction_readings DROP COLUMN created_at;

ALTER TABLE examination_findings DROP CONSTRAINT examination_findings_visit_id_key;
ALTER TABLE examination_findings DROP COLUMN updated_at;
ALTER TABLE examination_findings DROP COLUMN updated_by;
ALTER TABLE examination_findings DROP COLUMN bp;
ALTER TABLE examination_findings DROP COLUMN iop_time;
ALTER TABLE examination_findings DROP CONSTRAINT examination_findings_pkey;
ALTER TABLE examination_findings DROP COLUMN id;

ALTER TABLE visit_symptoms ALTER COLUMN double_vision DROP DEFAULT;
ALTER TABLE visit_symptoms ALTER COLUMN double_vision TYPE BOOLEAN USING (double_vision <> 'none');
ALTER TABLE visit_symptoms ALTER COLUMN double_vision SET DEFAULT FALSE;
ALTER TABLE visit_symptoms ALTER COLUMN double_vision DROP NOT NULL;

ALTER TABLE visit_symptoms ALTER COLUMN headache DROP DEFAULT;
ALTER TABLE visit_symptoms ALTER COLUMN headache TYPE BOOLEAN USING (headache <> 'none');
ALTER TABLE visit_symptoms ALTER COLUMN headache SET DEFAULT FALSE;
ALTER TABLE visit_symptoms ALTER COLUMN headache DROP NOT NULL;

ALTER TABLE visit_symptoms DROP CONSTRAINT visit_symptoms_visit_id_key;
