-- name: GetSourceByName :one
SELECT id, kind, name, licence FROM sources WHERE name = $1;
