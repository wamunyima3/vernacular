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
-- are decided one after another.
SELECT kind, language_code, status, author_id, open_for_review(status) AS open
FROM items WHERE id = $1 FOR UPDATE;
