-- +goose Up
ALTER TABLE clinic_settings ADD COLUMN workspace_slug VARCHAR(63) UNIQUE;
ALTER TABLE clinic_settings ADD COLUMN default_currency VARCHAR(3);
ALTER TABLE clinic_settings ADD COLUMN timezone VARCHAR(50);
ALTER TABLE clinic_settings ADD COLUMN require_consent_before_exam BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE clinic_settings ADD COLUMN low_stock_notifications BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE clinic_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    name VARCHAR(150) NOT NULL,
    address VARCHAR(255),
    exam_rooms INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'opening_soon')),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_clinic_locations_clinic_id ON clinic_locations (clinic_id);

-- One row per clinic for each of the following: clinic_id IS the primary
-- key, and GET/PUT upsert it the same way visit_symptoms upserts per visit.
CREATE TABLE clinic_notification_prefs (
    clinic_id UUID PRIMARY KEY REFERENCES clinic_settings(id),
    email_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    sms_appointment_reminders BOOLEAN NOT NULL DEFAULT FALSE,
    low_stock_alerts BOOLEAN NOT NULL DEFAULT TRUE,
    new_order_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    weekly_digest BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- What the clinic accepts from ITS patients at checkout — unrelated to the
-- clinic's own OpthalmOK subscription billing (billing/subscriptions tables).
-- bank_* and qr_* columns are only meaningful when the matching method flag
-- is true; nothing enforces that at the DB level, same as every other
-- nullable "only relevant when enabled" column in this schema.
CREATE TABLE clinic_payment_config (
    clinic_id UUID PRIMARY KEY REFERENCES clinic_settings(id),
    tax_label VARCHAR(20) NOT NULL DEFAULT 'VAT',
    tax_rate NUMERIC(5, 2) NOT NULL DEFAULT 0,
    vat_number VARCHAR(50),
    pan_number VARCHAR(50),
    method_cash BOOLEAN NOT NULL DEFAULT TRUE,
    method_card BOOLEAN NOT NULL DEFAULT FALSE,
    method_bank_transfer BOOLEAN NOT NULL DEFAULT FALSE,
    method_qr_wallet BOOLEAN NOT NULL DEFAULT FALSE,
    bank_name VARCHAR(150),
    bank_account_holder VARCHAR(150),
    bank_account_number VARCHAR(50),
    bank_branch VARCHAR(150),
    bank_swift_code VARCHAR(20),
    qr_provider VARCHAR(20)
        CHECK (qr_provider IS NULL OR qr_provider IN ('eSewa', 'Khalti', 'FonePay', 'other')),
    qr_merchant_id VARCHAR(100),
    qr_asset_url VARCHAR(500),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE clinic_pos_config (
    clinic_id UUID PRIMARY KEY REFERENCES clinic_settings(id),
    barcode_symbology VARCHAR(10) NOT NULL DEFAULT 'code128'
        CHECK (barcode_symbology IN ('code39', 'code128', 'ean13')),
    auto_generate_barcode BOOLEAN NOT NULL DEFAULT TRUE,
    barcode_scanning_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    auto_print_receipt BOOLEAN NOT NULL DEFAULT FALSE,
    require_pin_for_refunds BOOLEAN NOT NULL DEFAULT FALSE,
    receipt_printer VARCHAR(150),
    till_float_amount NUMERIC(10, 2),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE clinic_privacy_config (
    clinic_id UUID PRIMARY KEY REFERENCES clinic_settings(id),
    retention_years INTEGER NOT NULL DEFAULT 7,
    auto_anonymize_after_retention BOOLEAN NOT NULL DEFAULT FALSE,
    allow_patient_data_export BOOLEAN NOT NULL DEFAULT TRUE,
    allow_patient_deletion_requests BOOLEAN NOT NULL DEFAULT FALSE,
    require_2fa BOOLEAN NOT NULL DEFAULT FALSE,
    session_timeout_minutes INTEGER,
    login_alerts BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE clinic_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    actor_user_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    target_type VARCHAR(50),
    target_id VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_clinic_audit_log_clinic_id ON clinic_audit_log (clinic_id, created_at DESC);

-- access_token/refresh_token are written by the OAuth callback (not built
-- yet — /connect below is a stub) and are never selected by the repository
-- method the API read path uses, so they can't leak into a response.
CREATE TABLE clinic_integrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinic_settings(id),
    provider VARCHAR(30) NOT NULL
        CHECK (provider IN ('xero', 'twilio', 'google_calendar', 'essilor_connect')),
    connected BOOLEAN NOT NULL DEFAULT FALSE,
    connected_at TIMESTAMP,
    access_token VARCHAR(1000),
    refresh_token VARCHAR(1000),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (clinic_id, provider)
);

-- +goose Down
DROP TABLE clinic_integrations;
DROP TABLE clinic_audit_log;
DROP TABLE clinic_privacy_config;
DROP TABLE clinic_pos_config;
DROP TABLE clinic_payment_config;
DROP TABLE clinic_notification_prefs;
DROP TABLE clinic_locations;

ALTER TABLE clinic_settings DROP COLUMN low_stock_notifications;
ALTER TABLE clinic_settings DROP COLUMN require_consent_before_exam;
ALTER TABLE clinic_settings DROP COLUMN timezone;
ALTER TABLE clinic_settings DROP COLUMN default_currency;
ALTER TABLE clinic_settings DROP COLUMN workspace_slug;
