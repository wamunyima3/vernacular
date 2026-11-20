-- The review queue (P1-REV-02) asks for one language at a time. The agreement
-- view reused a CTE, so Postgres computed it for every language before
-- filtering. Window functions partitioned by language let a language filter
-- reach the base tables (Bemba: about 15ms instead of every language).
--
-- Priority becomes an enum in queue order, so min() picks the best.
-- open_for_review() is the one definition of which statuses still take reviews.

-- +goose Up
CREATE TYPE review_priority AS ENUM ('quick', 'standard', 'full');

-- +goose StatementBegin
CREATE FUNCTION open_for_review(status trust_status) RETURNS boolean
LANGUAGE sql IMMUTABLE AS $$ SELECT status IN ('machine', 'sourced', 'reviewed') $$;
-- +goose StatementEnd

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

-- +goose Down
DROP VIEW review_agreement;
DROP FUNCTION open_for_review(trust_status);
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
FROM per_sense s JOIN per_concept c USING (concept_id, language_code);
