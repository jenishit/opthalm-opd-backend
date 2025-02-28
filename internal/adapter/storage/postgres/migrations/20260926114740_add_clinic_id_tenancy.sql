-- +goose Up
-- This backend has no production data yet (everything in every environment
-- so far is test/dev data from this project's own development), so rather
-- than writing a real backfill strategy for a NOT NULL tenant column, we
-- clear the tables that are about to require one. A real backfill (assign
-- every existing row to some default clinic) would be the right move once
-- there's real data to preserve.
TRUNCATE TABLE
    lab_job_status_history, lab_jobs,
    stock_purchase_items, stock_purchases,
    inventory_stock_movements, inventory_items, vendors,
    payments, invoice_items, invoices,
    follow_ups, refraction_readings, investigations, examination_findings, visit_symptoms,
    visits,
    patients,
    profile, users
    CASCADE;

ALTER TABLE users ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_users_clinic_id ON users (clinic_id);

ALTER TABLE patients ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_patients_clinic_id ON patients (clinic_id);

ALTER TABLE visits ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_visits_clinic_id ON visits (clinic_id);

ALTER TABLE visit_symptoms ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE examination_findings ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE investigations ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE refraction_readings ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE follow_ups ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);

ALTER TABLE invoices ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_invoices_clinic_id ON invoices (clinic_id);
ALTER TABLE invoice_items ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE payments ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);

ALTER TABLE inventory_items ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_inventory_items_clinic_id ON inventory_items (clinic_id);
ALTER TABLE inventory_stock_movements ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE vendors ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_vendors_clinic_id ON vendors (clinic_id);
ALTER TABLE stock_purchases ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
ALTER TABLE stock_purchase_items ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);

ALTER TABLE lab_jobs ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);
CREATE INDEX idx_lab_jobs_clinic_id ON lab_jobs (clinic_id);
ALTER TABLE lab_job_status_history ADD COLUMN clinic_id UUID NOT NULL REFERENCES clinic_settings(id);

CREATE TABLE subscriptions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id          UUID NOT NULL UNIQUE REFERENCES clinic_settings(id),
    plan_name          VARCHAR(50) NOT NULL DEFAULT 'trial',
    status             VARCHAR(20) NOT NULL DEFAULT 'trialing'
                           CHECK (status IN ('trialing', 'active', 'past_due', 'cancelled')),
    current_period_end TIMESTAMP NOT NULL,
