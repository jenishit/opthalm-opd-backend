-- +goose Up
CREATE TABLE visit_symptoms (
    visit_id UUID REFERENCES visits(id) NOT NULL,
    foreign_body_sensation VARCHAR(4) NOT NULL DEFAULT 'none',
    redness VARCHAR(4) NOT NULL DEFAULT 'none',
    swelling VARCHAR(4) NOT NULL DEFAULT 'none',
    diminished_distance_vision VARCHAR(4) NOT NULL DEFAULT 'none',
    diminished_near_vision VARCHAR(4) NOT NULL DEFAULT 'none',
    blurring VARCHAR(4) NOT NULL DEFAULT 'none',
    headache BOOLEAN DEFAULT FALSE,
    eye_pain VARCHAR(4) NOT NULL DEFAULT 'none',
    watery_eyes VARCHAR(4) NOT NULL DEFAULT 'none',
    discharge VARCHAR(4) NOT NULL DEFAULT 'none',
    itching VARCHAR(4) NOT NULL DEFAULT 'none',
    photophobia VARCHAR(4) NOT NULL DEFAULT 'none',
    floaters VARCHAR(4) NOT NULL DEFAULT 'none',
    flashes VARCHAR(4) NOT NULL DEFAULT 'none',
    double_vision BOOLEAN DEFAULT FALSE,
    recent_fever BOOLEAN DEFAULT FALSE,
    complaint_duration VARCHAR(50),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE examination_findings (
    visit_id UUID REFERENCES visits(id),
    iop_od NUMERIC(4, 1),
    iop_os NUMERIC(4, 1),
    convergence VARCHAR(100),
    cover_test VARCHAR(100),
    eom VARCHAR(100),
    anterior_segment_od VARCHAR(255),
    anterior_segment_os VARCHAR(255),
    fundus_od VARCHAR(255),
    fundus_os VARCHAR(255),
    created_by UUID REFERENCES users(id)
);

CREATE TABLE investigations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID REFERENCES visits(id) NOT NULL,
    type VARCHAR(50) NOT NULL,
    eye VARCHAR(2),
    result_value VARCHAR(50),
    notes VARCHAR(255),
    created_by UUID REFERENCES users(id)
);

CREATE TABLE refraction_readings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID REFERENCES visits(id) NOT NULL,
    eye VARCHAR(2) NOT NULL,
    stage VARCHAR(25) NOT NULL,
