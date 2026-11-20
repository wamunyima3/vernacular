-- name: GetSourceByName :one
SELECT id, kind, name, licence FROM sources WHERE name = $1;

-- name: ListSources :many
SELECT * FROM sources ORDER BY name;

-- name: GetSource :one
SELECT * FROM sources WHERE id = $1;

-- name: CreateSource :one
INSERT INTO sources (kind, name, rights_holder, licence, url, allowed_uses, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateSource :one
UPDATE sources SET kind = $2, name = $3, rights_holder = $4, licence = $5, url = $6,
    allowed_uses = $7, notes = $8
WHERE id = $1
RETURNING *;

-- name: DeleteSource :execrows
DELETE FROM sources WHERE id = $1;

-- name: SetPermissionKey :one
UPDATE sources SET permission_key = $2 WHERE id = $1 RETURNING *;
