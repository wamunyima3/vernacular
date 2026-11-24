-- Phrase candidates join the review queue (P1-AI-04b). The agreement view
-- counts phrases as it counts senses, by the concept they express, and gains
-- a kind, so a card holds the senses or the phrases of one concept, never
-- both: agreement and priority are worked out within a kind.

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
        ) attested
    ) ranked
) counted
WHERE n = 1;

-- +goose Down
DROP VIEW review_agreement;
CREATE VIEW review_agreement AS
SELECT item_id, language_code, concept_id, agreement, answered,
       (CASE
           WHEN agreement >= 2 AND agreement = answered THEN 'quick'
           WHEN agreement >= 2 THEN 'standard'
           ELSE 'full'
       END)::review_priority AS priority
FROM (
    SELECT item_id, language_code, concept_id, agreement, n,
           max(source_rank) OVER (PARTITION BY language_code, concept_id)::int AS answered
    FROM (
        SELECT item_id, language_code, concept_id,
               count(*) OVER (PARTITION BY language_code, item_id)::int AS agreement,
               dense_rank() OVER (PARTITION BY language_code, concept_id ORDER BY source_id) AS source_rank,
               row_number() OVER (PARTITION BY language_code, item_id) AS n
        FROM (
            SELECT DISTINCT s.id AS item_id, e.language_code, s.concept_id, isrc.source_id
            FROM senses s
            JOIN entries e ON e.id = s.entry_id
            JOIN item_sources isrc ON isrc.item_id = s.id
            WHERE s.concept_id IS NOT NULL
        ) attested
    ) ranked
) counted
WHERE n = 1;
