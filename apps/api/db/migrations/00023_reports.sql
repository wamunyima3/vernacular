-- Reports of wrong items (PRV-3, DQ-8, US-15, P1-API-10). Anyone can report
-- a sense; a reviewer of its language confirms or dismisses the report. A
-- confirmed report makes the item disputed (review.Status) until a senior
-- reviewer's resolution settles it, which sets settled_by, as a send-back
-- sets reviews aside.

-- +goose Up
CREATE TYPE report_state AS ENUM ('open', 'confirmed', 'dismissed');

CREATE TABLE reports (
    id          uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id     uuid         NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    reason      text         NOT NULL CHECK (reason <> '' AND length(reason) <= 1000),
    reporter_id text         REFERENCES users (id) ON DELETE SET NULL, -- null when signed out
    state       report_state NOT NULL DEFAULT 'open',
    handled_by  text         REFERENCES users (id) ON DELETE SET NULL,
    handled_at  timestamptz,
    note        text         NOT NULL DEFAULT '',
    settled_by  uuid         REFERENCES resolutions (id) ON DELETE SET NULL,
    created_at  timestamptz  NOT NULL DEFAULT now(),
    updated_at  timestamptz  NOT NULL DEFAULT now(),
    CHECK ((state = 'open') = (handled_at IS NULL))
);
CREATE INDEX reports_open_idx ON reports (created_at) WHERE state = 'open';
CREATE INDEX reports_item_idx ON reports (item_id);
SELECT track('reports');

-- +goose Down
DROP TABLE reports;
DROP TYPE report_state;
