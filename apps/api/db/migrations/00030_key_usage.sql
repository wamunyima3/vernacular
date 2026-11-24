-- Calls per API key per day (DEV-5, US-13), for the usage dashboard. The API
-- counts in memory and adds its counts here every few seconds, so a busy key
-- costs one write per flush, not one per call. Days are UTC.

-- +goose Up
CREATE TABLE key_usage (
    key_id  uuid NOT NULL REFERENCES api_keys (id) ON DELETE CASCADE,
    day     date NOT NULL,
    calls   int  NOT NULL DEFAULT 0 CHECK (calls >= 0),
    errors  int  NOT NULL DEFAULT 0 CHECK (errors >= 0),  -- 4xx and 5xx, other than 429
    limited int  NOT NULL DEFAULT 0 CHECK (limited >= 0), -- 429: over the rate limit
    PRIMARY KEY (key_id, day)
);

-- +goose Down
DROP TABLE key_usage;
