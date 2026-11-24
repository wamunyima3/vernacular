-- Example sentences go through the review queue like meanings and phrases
-- (LEX-5, P2-09a): a card holds the examples for one concept in a language,
-- and two accepts verify one.

-- +goose Up
CREATE OR REPLACE VIEW review_agreement AS
SELECT item_id, language_code, concept_id, agreement, answered,
       (CASE
           WHEN agreement >= 2 AND agreement = answered THEN 'quick'
           WHEN agreement >= 2 THEN 'standard'
           ELSE 'full'
       END)::review_priority AS priority,
       kind
FROM (
    SELECT item_id, language_code, concept_id, kind, agreement, n,
           max(source_rank) OVER (PARTITION BY language_code, concept_id, kind)::int AS answered
    FROM (
        SELECT item_id, language_code, concept_id, kind,
               count(*) OVER (PARTITION BY language_code, item_id)::int AS agreement,
               dense_rank() OVER (PARTITION BY language_code, concept_id, kind ORDER BY source_id) AS source_rank,
               row_number() OVER (PARTITION BY language_code, item_id) AS n
        FROM (
            SELECT DISTINCT s.id AS item_id, e.language_code, s.concept_id, 'sense'::item_kind AS kind, isrc.source_id
            FROM senses s
            JOIN entries e ON e.id = s.entry_id
            JOIN item_sources isrc ON isrc.item_id = s.id
            WHERE s.concept_id IS NOT NULL
            UNION ALL -- the kinds differ, so no row is in both
            SELECT DISTINCT p.id, p.language_code, p.concept_id, 'phrase'::item_kind, isrc.source_id
            FROM phrases p
            JOIN item_sources isrc ON isrc.item_id = p.id
            WHERE p.concept_id IS NOT NULL
            UNION ALL
            SELECT DISTINCT x.id, e.language_code, s.concept_id, 'example'::item_kind, isrc.source_id
            FROM examples x
            JOIN senses s ON s.id = x.sense_id
            JOIN entries e ON e.id = s.entry_id
            JOIN item_sources isrc ON isrc.item_id = x.id
            WHERE s.concept_id IS NOT NULL
        ) attested
    ) ranked
) counted
WHERE n = 1;

-- +goose Down
CREATE OR REPLACE VIEW review_agreement AS
SELECT item_id, language_code, concept_id, agreement, answered,
       (CASE
           WHEN agreement >= 2 AND agreement = answered THEN 'quick'
           WHEN agreement >= 2 THEN 'standard'
           ELSE 'full'
       END)::review_priority AS priority,
       kind
FROM (
    SELECT item_id, language_code, concept_id, kind, agreement, n,
           max(source_rank) OVER (PARTITION BY language_code, concept_id, kind)::int AS answered
    FROM (
        SELECT item_id, language_code, concept_id, kind,
               count(*) OVER (PARTITION BY language_code, item_id)::int AS agreement,
               dense_rank() OVER (PARTITION BY language_code, concept_id, kind ORDER BY source_id) AS source_rank,
               row_number() OVER (PARTITION BY language_code, item_id) AS n
        FROM (
            SELECT DISTINCT s.id AS item_id, e.language_code, s.concept_id, 'sense'::item_kind AS kind, isrc.source_id
            FROM senses s
            JOIN entries e ON e.id = s.entry_id
            JOIN item_sources isrc ON isrc.item_id = s.id
            WHERE s.concept_id IS NOT NULL
            UNION ALL -- the kinds differ, so no row is in both
            SELECT DISTINCT p.id, p.language_code, p.concept_id, 'phrase'::item_kind, isrc.source_id
            FROM phrases p
            JOIN item_sources isrc ON isrc.item_id = p.id
            WHERE p.concept_id IS NOT NULL
        ) attested
    ) ranked
) counted
WHERE n = 1;
