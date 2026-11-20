-- name: ListLanguages :many
SELECT code, name, own_name, wave, status
FROM languages
ORDER BY wave, name;

-- name: GetLanguageStatus :one
SELECT name, status FROM languages WHERE code = $1;
