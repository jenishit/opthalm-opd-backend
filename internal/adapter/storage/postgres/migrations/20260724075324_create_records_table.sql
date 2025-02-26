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
