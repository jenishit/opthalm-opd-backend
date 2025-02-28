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
