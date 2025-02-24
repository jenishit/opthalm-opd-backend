-- +goose Up
CREATE TABLE clinic_settings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lab_name        VARCHAR(150) NOT NULL,
    tagline         VARCHAR(255),
    address         VARCHAR(255),
    phone           VARCHAR(20) UNIQUE,
    email           VARCHAR(150) UNIQUE,
    registration_no VARCHAR(50) UNIQUE,
    report_footer   TEXT,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by      UUID REFERENCES users(id)
);

-- +goose Down
DROP TABLE clinic_settings;
