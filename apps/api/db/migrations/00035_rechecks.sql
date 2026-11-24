-- Quarterly re-checks (DQ-6, US-28). Each quarter 5% of each language's
-- verified items are drawn and checked again by reviewers who neither wrote
-- nor reviewed them. Each sample keeps the batch the item came from, its
-- first source and model, so a batch failing more than 2% can be sent back.

-- +goose Up
CREATE TYPE recheck_outcome AS ENUM ('pass', 'fail');

CREATE TABLE rechecks (
    quarter       text            NOT NULL CHECK (quarter ~ '^[0-9]{4}-Q[1-4]$'),
    item_id       uuid            NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    language_code text            NOT NULL REFERENCES languages (code),
    source_id     uuid            NOT NULL REFERENCES sources (id),
    model         text            NOT NULL DEFAULT '',
    outcome       recheck_outcome,
    note          text,
    checked_by    text,
    checked_at    timestamptz,
    created_at    timestamptz     NOT NULL DEFAULT now(),
    updated_at    timestamptz     NOT NULL DEFAULT now(),
    PRIMARY KEY (quarter, item_id),
    CHECK ((outcome IS NULL) = (checked_at IS NULL) AND (outcome IS NULL) = (checked_by IS NULL)),
    CHECK (outcome IS DISTINCT FROM 'fail' OR note <> '')
);
CREATE INDEX rechecks_open_idx ON rechecks (language_code, created_at) WHERE outcome IS NULL;
SELECT track('rechecks');

-- +goose Down
DROP TABLE rechecks;
DROP TYPE recheck_outcome;
