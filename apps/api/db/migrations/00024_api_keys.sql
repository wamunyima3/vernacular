-- Developer API keys (DEV-1, NFR-5, US-01, P1-API-11). A key is shown once
-- when created; only its SHA-256 hash is stored, with a short prefix so the
-- owner can tell keys apart. Revoked keys stay for the audit trail.

-- +goose Up
CREATE TABLE api_keys (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      text        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name         text        NOT NULL CHECK (name <> '' AND length(name) <= 100),
    prefix       text        NOT NULL,
    hash         bytea       NOT NULL UNIQUE,
    last_used_at timestamptz,
    revoked_at   timestamptz,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX api_keys_user_idx ON api_keys (user_id, created_at DESC);
SELECT track('api_keys');

-- +goose Down
DROP TABLE api_keys;
