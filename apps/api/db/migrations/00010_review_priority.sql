-- Review priority from agreement between sources (P1-AI-03, AI-3, US-21).
-- agreement is how many distinct sources attest the item; the queue
-- (P1-REV-02) serves quick confirms first.

-- +goose Up
CREATE TYPE review_priority AS ENUM ('quick', 'standard', 'full');
ALTER TABLE items
    ADD COLUMN review_priority review_priority NOT NULL DEFAULT 'full',
    ADD COLUMN agreement       smallint        NOT NULL DEFAULT 1;
CREATE INDEX items_priority_idx ON items (language_code, status, review_priority);

-- +goose Down
ALTER TABLE items DROP COLUMN review_priority, DROP COLUMN agreement;
DROP TYPE review_priority;
