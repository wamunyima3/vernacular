-- name: UpsertConcept :one
-- Concepts are shared across languages. A later import can promote one to the
-- core list, or give a 'general' concept its domain, but never demotes it.
-- Rows that would not change are not written, so re-runs add no audit rows.
WITH written AS (
    INSERT INTO concepts (gloss_en, core, domain) VALUES (@gloss_en, @core, @domain)
    ON CONFLICT (gloss_en, definition) DO UPDATE SET
        core   = concepts.core OR excluded.core,
        domain = CASE WHEN concepts.domain = 'general' THEN excluded.domain ELSE concepts.domain END
    WHERE (excluded.core AND NOT concepts.core)
       OR (concepts.domain = 'general' AND excluded.domain <> 'general')
    RETURNING id, (xmax = 0) AS inserted
)
SELECT id, inserted FROM written
UNION ALL
SELECT id, false FROM concepts
WHERE gloss_en = @gloss_en AND definition = '' AND NOT EXISTS (SELECT 1 FROM written);

-- name: FindEntry :one
SELECT id FROM entries
WHERE language_code = $1 AND headword = $2 AND part_of_speech IS NOT DISTINCT FROM $3;

-- name: CreateEntry :exec
INSERT INTO entries (id, language_code, headword, part_of_speech) VALUES ($1, $2, $3, $4);

-- name: FindSense :one
SELECT id FROM senses WHERE entry_id = $1 AND concept_id = $2;

-- name: CreateSense :exec
INSERT INTO senses (id, entry_id, concept_id, gloss, position)
VALUES ($1, $2, $3, $4, (SELECT coalesce(max(position), 0) + 1 FROM senses WHERE entry_id = $2));
