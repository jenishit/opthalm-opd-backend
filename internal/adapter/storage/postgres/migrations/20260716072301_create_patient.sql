-- +goose Up
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address VARCHAR(255),
    dob DATE NOT NULL,
    gender VARCHAR(10) NOT NULL DEFAULT 'other'
                     CHECK (gender IN ('male','female','other')),
    occupation VARCHAR(100),
    registered_on DATE NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    deleted_at TIMESTAMP DEFAULT NULL,
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE patients;
