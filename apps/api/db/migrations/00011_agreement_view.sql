-- Review priority is derived, not stored (P1-AI-03a). Agreement is how many
-- distinct sources attest a sense; any write that adds a source changes it,
-- so storing it meant re-scoring on every path. The review queue (P1-REV-02)
-- reads this view instead.
--
-- Priority per sense of a concept in a language (AI-3, US-21):
--   quick    every source that answered for the concept gives this word (2+)
--   standard at least two sources give it
--   full     otherwise

-- +goose Up
DROP INDEX items_priority_idx;
ALTER TABLE items DROP COLUMN review_priority, DROP COLUMN agreement;
DROP TYPE review_priority;

CREATE VIEW review_agreement AS
WITH attested AS (
    SELECT s.id AS item_id, e.language_code, s.concept_id, isrc.source_id
    FROM senses s
    JOIN entries e ON e.id = s.entry_id
    JOIN item_sources isrc ON isrc.item_id = s.id
    WHERE s.concept_id IS NOT NULL
    GROUP BY s.id, e.language_code, s.concept_id, isrc.source_id
), per_sense AS (
    SELECT item_id, language_code, concept_id, count(*)::int AS agreement
    FROM attested GROUP BY item_id, language_code, concept_id
), per_concept AS (
    SELECT language_code, concept_id, count(DISTINCT source_id)::int AS answered
    FROM attested GROUP BY language_code, concept_id
)
SELECT s.item_id, s.language_code, s.concept_id, s.agreement, c.answered,
       CASE
           WHEN s.agreement >= 2 AND s.agreement = c.answered THEN 'quick'
           WHEN s.agreement >= 2 THEN 'standard'
           ELSE 'full'
       END AS priority
FROM per_sense s JOIN per_concept c USING (language_code, concept_id);

-- +goose Down
DROP VIEW review_agreement;
CREATE TYPE review_priority AS ENUM ('quick', 'standard', 'full');
ALTER TABLE items
    ADD COLUMN review_priority review_priority NOT NULL DEFAULT 'full',
    ADD COLUMN agreement       smallint        NOT NULL DEFAULT 1;
CREATE INDEX items_priority_idx ON items (language_code, status, review_priority);
