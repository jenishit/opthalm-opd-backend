-- +goose Up
CREATE TABLE sessions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL REFERENCES users(id),
    refresh_token_hash TEXT NOT NULL,
    user_agent         VARCHAR(255),
    ip_address         VARCHAR(45),
    issued_at          TIMESTAMP NOT NULL DEFAULT now(),
    expires_at         TIMESTAMP NOT NULL,
    revoked_at         TIMESTAMP,
    replaced_by        UUID REFERENCES sessions(id)
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE UNIQUE INDEX idx_sessions_refresh_token_hash ON sessions (refresh_token_hash);

-- +goose Down
DROP TABLE sessions;
