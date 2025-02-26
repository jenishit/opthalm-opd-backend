-- +goose Up
CREATE TABLE history_conditions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    deleted_at  TIMESTAMP DEFAULT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE diagnosis_catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    icd10_code VARCHAR(10) UNIQUE,
    name VARCHAR(255) NOT NULL,
    deleted_at  TIMESTAMP DEFAULT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE medicines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medicine_name VARCHAR(150) NOT NULL,
    brand_name VARCHAR(150),
    strength VARCHAR(50),
    form VARCHAR(50),
    deleted_at  TIMESTAMP DEFAULT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE medicines;
DROP TABLE diagnosis_catalog;
DROP TABLE history_conditions;
