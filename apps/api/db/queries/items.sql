-- name: CreateItem :one
INSERT INTO items (kind, language_code, status, author_id) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: LinkSource :exec
INSERT INTO item_sources (item_id, source_id, reference, model) VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;

-- name: ListItemSources :many
SELECT DISTINCT item_id, source_id FROM item_sources WHERE item_id = ANY(@ids::uuid[]);

-- name: SetPriority :exec
-- Unchanged rows are not written, so re-scoring adds no audit rows.
UPDATE items SET review_priority = @review_priority, agreement = @agreement
WHERE id = @id AND (review_priority, agreement) IS DISTINCT FROM (@review_priority::review_priority, @agreement::smallint);
