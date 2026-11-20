-- name: CreateItem :one
INSERT INTO items (kind, language_code, status, author_id) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: LinkSource :exec
INSERT INTO item_sources (item_id, source_id, reference, model) VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;
