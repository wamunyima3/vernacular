-- Senior reviewers settle disputes (P1-REV-04, CON-4). A resolution is an
-- input to the promotion rule: `keep` verifies the item, `send_back` returns
-- it to the queue and sets its reviews aside, so only new ones count. Set-aside
-- reviews stay for history and still stop their reviewer reviewing it again.

-- +goose Up
CREATE TYPE resolution_outcome AS ENUM ('keep', 'send_back');

CREATE TABLE resolutions (
    id         uuid               PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id    uuid               NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    senior_id  text               NOT NULL,
    outcome    resolution_outcome NOT NULL,
    note       text               NOT NULL CHECK (note <> ''),
    created_at timestamptz        NOT NULL DEFAULT now(),
    updated_at timestamptz        NOT NULL DEFAULT now()
);
CREATE INDEX resolutions_item_idx ON resolutions (item_id, created_at DESC);
SELECT track('resolutions');

ALTER TABLE reviews ADD COLUMN set_aside_by uuid REFERENCES resolutions (id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE reviews DROP COLUMN set_aside_by;
DROP TABLE resolutions;
DROP TYPE resolution_outcome;
