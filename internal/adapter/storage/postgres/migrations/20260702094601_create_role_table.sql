-- +goose Up
CREATE TABLE role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(50) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS "role";
