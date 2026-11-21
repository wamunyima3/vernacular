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

-- name: ListCoreConcepts :many
SELECT id FROM concepts WHERE core ORDER BY gloss_en;

-- name: LockLanguage :exec
-- Serialises lexicon writes per language across concurrent transactions;
-- released when the transaction ends.
SELECT pg_advisory_xact_lock(hashtext('lexicon:' || @language_code::text));

-- name: CreateExample :exec
INSERT INTO examples (id, sense_id, text, translation_en) VALUES ($1, $2, $3, $4);

-- name: ListSenseCards :many
SELECT s.id, s.concept_id, c.gloss_en, e.headword, e.part_of_speech, e.noun_class, e.plural
FROM senses s JOIN entries e ON e.id = s.entry_id
LEFT JOIN concepts c ON c.id = s.concept_id
WHERE s.id = ANY(@ids::uuid[])
ORDER BY array_position(@ids::uuid[], s.id);

-- name: CountEntries :many
-- Per language: entries with a verified sense, and entries with any sense that
-- can be served (US-03, P1-API-01).
SELECT e.language_code,
       count(DISTINCT e.id) FILTER (WHERE i.status = 'verified')::int AS verified,
       count(DISTINCT e.id)::int AS total
FROM entries e
JOIN senses s ON s.entry_id = e.id
JOIN items i ON i.id = s.id
WHERE servable(i.status)
GROUP BY e.language_code;

-- name: SearchEntries :many
-- Entry search (LEX-4, US-05). Each entry keeps its best match: the headword,
-- a spelling variant (stored, or the style guide's spelling of the query), the
-- search stem, or a similar spelling by trigram. Only servable entries of the variety (any when
-- none) with a sense at or above min_status are returned.
WITH hits AS (
    SELECT e.id, 'headword' AS matched, 1 AS rank, 1::real AS score
    FROM entries e WHERE e.language_code = @language AND e.headword_fold = @fold::text
    UNION ALL
    SELECT e.id, 'variant', 2, 1 FROM entries e
    WHERE e.language_code = @language AND e.headword_fold = @standard::text AND @standard::text <> @fold::text
    UNION ALL
    SELECT v.entry_id, 'variant', 2, 1 FROM variants v JOIN entries e ON e.id = v.entry_id
    WHERE e.language_code = @language AND v.fold = @fold::text
    UNION ALL
    SELECT e.id, 'root', 3, 1 FROM entries e WHERE e.language_code = @language AND e.stem = @stem::text
    UNION ALL
    SELECT e.id, 'similar', 4, similarity(e.headword_fold, @fold::text)
    FROM entries e WHERE e.language_code = @language AND e.headword_fold % @fold::text
), best AS (
    SELECT DISTINCT ON (id) id, matched, rank, score FROM hits ORDER BY id, rank, score DESC
)
SELECT e.id, b.matched::text AS matched
FROM best b
JOIN entries e ON e.id = b.id
JOIN items i ON i.id = e.id
WHERE servable(i.status)
  AND (sqlc.narg(variety)::text IS NULL OR e.variety = sqlc.narg(variety)::text)
  AND EXISTS (SELECT 1 FROM senses s JOIN items si ON si.id = s.id
              WHERE s.entry_id = e.id AND at_least(si.status, @min_status))
ORDER BY b.rank, b.score DESC, e.headword, e.id
LIMIT @lim;

-- name: ListSenses :many
-- The senses of entries at or above min_status, in order.
SELECT s.id, s.entry_id, s.gloss, s.definition, s.concept_id
FROM senses s JOIN items i ON i.id = s.id
WHERE s.entry_id = ANY(@entry_ids::uuid[]) AND at_least(i.status, @min_status)
ORDER BY s.entry_id, s.position, s.id;

-- name: ListEntries :many
-- Servable entries by ID (US-16); the caller keeps its own order.
SELECT e.id, e.language_code, e.headword, e.part_of_speech, e.noun_class, e.plural, e.root,
       e.variety, e.register
FROM entries e JOIN items i ON i.id = e.id
WHERE e.id = ANY(@ids::uuid[]) AND servable(i.status);

-- name: ListVariants :many
SELECT spelling, note FROM variants WHERE entry_id = $1 ORDER BY spelling;

-- name: ListExamples :many
SELECT x.id, x.sense_id, x.text, x.translation_en
FROM examples x JOIN items i ON i.id = x.id
WHERE x.sense_id = ANY(@sense_ids::uuid[]) AND at_least(i.status, @min_status)
ORDER BY x.sense_id, x.created_at, x.id;

-- name: GetConcept :one
SELECT id, gloss_en, definition, domain, core FROM concepts WHERE id = $1;

-- name: ListConceptEntries :many
-- Entries with a sense at or above min_status expressing one of the concepts
-- (LEX-7), in the given languages and variety, or in all when none are given. At most 200 per concept.
WITH found AS (
    SELECT DISTINCT s.concept_id, e.id, e.language_code, e.headword
    FROM senses s
    JOIN items i ON i.id = s.id
    JOIN entries e ON e.id = s.entry_id
    WHERE s.concept_id = ANY(@concept_ids::uuid[]) AND at_least(i.status, @min_status)
      AND (sqlc.narg(variety)::text IS NULL OR e.variety = sqlc.narg(variety)::text)
      AND (coalesce(cardinality(@languages::text[]), 0) = 0 OR e.language_code = ANY(@languages::text[]))
), ranked AS (
    SELECT *, row_number() OVER (PARTITION BY concept_id ORDER BY language_code, headword, id) AS n
    FROM found
)
SELECT concept_id, id FROM ranked WHERE n <= 200
ORDER BY language_code, headword, id;

-- name: FindConcepts :many
-- Concepts by English gloss (LEX-2): exact first, then similar by trigram.
SELECT id, gloss_en, fold(gloss_en) = @fold::text AS exact
FROM concepts
WHERE fold(gloss_en) = @fold::text OR fold(gloss_en) % @fold::text
ORDER BY exact DESC, similarity(fold(gloss_en), @fold::text) DESC, gloss_en, id
LIMIT @lim;

-- name: ListConcepts :many
SELECT id, gloss_en FROM concepts WHERE id = ANY(@ids::uuid[]);
