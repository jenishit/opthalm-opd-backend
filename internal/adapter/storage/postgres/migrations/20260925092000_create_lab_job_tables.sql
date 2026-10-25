-- +goose Up
CREATE TABLE lab_jobs (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id             UUID REFERENCES invoices(id),
    invoice_item_id        UUID REFERENCES invoice_items(id),
    patient_id             UUID NOT NULL REFERENCES patients(id),
    vendor_id              UUID REFERENCES vendors(id),
    job_type               VARCHAR(50) NOT NULL,
    status                 VARCHAR(20) NOT NULL DEFAULT 'in_fitting'
                               CHECK (status IN ('in_fitting', 'ready_to_deliver', 'delivered', 'cancelled')),
    expected_delivery_date DATE,
    delivered_at           TIMESTAMP,
    advance_payment        NUMERIC(10, 2) NOT NULL DEFAULT 0,
    notes                  VARCHAR(255),
    created_by             UUID NOT NULL REFERENCES users(id),
    updated_by             UUID NOT NULL REFERENCES users(id),
    created_at             TIMESTAMP NOT NULL DEFAULT now(),
    updated_at             TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_lab_jobs_patient_id ON lab_jobs (patient_id);

CREATE TABLE lab_job_status_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lab_job_id  UUID NOT NULL REFERENCES lab_jobs(id),
    status      VARCHAR(20) NOT NULL,
    changed_at  TIMESTAMP NOT NULL DEFAULT now(),
    changed_by  UUID NOT NULL REFERENCES users(id),
    notes       VARCHAR(255)
);

CREATE INDEX idx_lab_job_status_history_job_id ON lab_job_status_history (lab_job_id);

-- +goose Down
DROP TABLE lab_job_status_history;
DROP TABLE lab_jobs;
