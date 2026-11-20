-- name: ListLanguages :many
SELECT code, name, own_name, wave, status
FROM languages
ORDER BY wave, name;
