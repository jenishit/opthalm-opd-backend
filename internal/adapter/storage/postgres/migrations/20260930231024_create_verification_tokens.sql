-- +goose Up
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP;

CREATE TABLE verification_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id),
    purpose    VARCHAR(30) NOT NULL CHECK (purpose IN ('password_reset', 'email_verification')),
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at    TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_verification_tokens_user_purpose ON verification_tokens (user_id, purpose);
CREATE UNIQUE INDEX idx_verification_tokens_token_hash ON verification_tokens (token_hash);

-- +goose Down
DROP TABLE verification_tokens;
ALTER TABLE users DROP COLUMN email_verified_at;
