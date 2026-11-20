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
INSERT INTO entries (id, language_code, headword, part_of_speech, stem, noun_class, plural)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindSense :one
SELECT id FROM senses WHERE entry_id = $1 AND concept_id = $2;

-- name: CreateSense :exec
INSERT INTO senses (id, entry_id, concept_id, gloss, position)
VALUES ($1, $2, $3, $4, (SELECT coalesce(max(position), 0) + 1 FROM senses WHERE entry_id = $2));

-- name: AddVariant :execrows
INSERT INTO variants (entry_id, spelling, note) VALUES ($1, $2, $3)
ON CONFLICT (entry_id, spelling) DO NOTHING;

-- name: ListHeadwords :many
SELECT id, headword, stem FROM entries WHERE language_code = $1;

-- name: SetStem :exec
UPDATE entries SET stem = $2 WHERE id = $1;

-- name: GetConcept :one
SELECT id, gloss_en, definition, domain FROM concepts WHERE id = $1;

-- name: ListCoreConcepts :many
SELECT id FROM concepts WHERE core ORDER BY gloss_en;

-- name: LockLanguage :exec
-- Serialises lexicon writes per language across concurrent transactions;
-- released when the transaction ends.
SELECT pg_advisory_xact_lock(hashtext('lexicon:' || @language_code::text));

-- name: ListSensesOfConcept :many
SELECT s.id FROM senses s JOIN entries e ON e.id = s.entry_id
WHERE e.language_code = $1 AND s.concept_id = $2;

-- name: ListConceptsWithSenses :many
SELECT DISTINCT s.concept_id FROM senses s JOIN entries e ON e.id = s.entry_id
WHERE e.language_code = $1 AND s.concept_id IS NOT NULL;
