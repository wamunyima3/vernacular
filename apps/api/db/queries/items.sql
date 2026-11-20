-- name: CreateItem :one
INSERT INTO items (kind, language_code, status, author_id) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: LinkSource :exec
INSERT INTO item_sources (item_id, source_id, reference, model) VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;


-- name: GetTrust :one
-- The trust fields every returned item carries (PRV-1).
SELECT i.status,
       count(r.id) FILTER (WHERE r.decision = 'accept')::int  AS accepts,
       count(r.id) FILTER (WHERE r.decision = 'correct')::int AS corrections
FROM items i LEFT JOIN reviews r ON r.item_id = i.id
WHERE i.id = $1
GROUP BY i.id;

-- name: ListTrustSources :many
SELECT s.name, s.licence, s.url, isrc.reference, isrc.model
FROM item_sources isrc JOIN sources s ON s.id = isrc.source_id
WHERE isrc.item_id = $1
ORDER BY s.name, isrc.reference;

-- name: LockItem :one
-- Locks the item until the transaction ends, so concurrent reviews of it
-- are decided one after another. NO KEY UPDATE, because only status changes:
-- writers linking sources to the item (which hold the lexicon's language
-- lock) are not blocked, so they cannot deadlock with a review.
-- origin is the status the item had before any review: `machine` when only
-- machine sources cite it, otherwise `sourced` (ADR-0004).
SELECT i.kind, i.language_code, i.status, i.author_id, open_for_review(i.status) AS open,
       (SELECT (CASE WHEN bool_or(s.kind <> 'machine') THEN 'sourced' ELSE 'machine' END)::trust_status
        FROM item_sources isrc JOIN sources s ON s.id = isrc.source_id
        WHERE isrc.item_id = i.id) AS origin
FROM items i WHERE i.id = $1 FOR NO KEY UPDATE OF i;


-- name: SetItemStatus :exec
-- Only the review module's promotion function calls this (AGENTS.md rule 2).
UPDATE items SET status = $2, updated_at = now() WHERE id = $1;
