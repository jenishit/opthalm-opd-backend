-- +goose Up
CREATE TABLE patient_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    uses_spectacles BOOLEAN NOT NULL DEFAULT FALSE,
    systemic_medications VARCHAR(500),
    topical_medications VARCHAR(500),
    notes VARCHAR(1000),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_patient_history_patient_id ON patient_history (patient_id);

-- A history entry can reference several catalog conditions at once (e.g.
-- "diabetic" + "hypertensive" recorded in the same note), hence a join
-- table rather than a column on patient_history.
CREATE TABLE patient_history_conditions (
    history_id UUID NOT NULL REFERENCES patient_history(id),
    condition_id UUID NOT NULL REFERENCES history_conditions(id),
    PRIMARY KEY (history_id, condition_id)
);

-- +goose Down
DROP TABLE patient_history_conditions;
DROP TABLE patient_history;
