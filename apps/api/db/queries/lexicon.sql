-- name: UpsertConcept :one
-- Concepts are shared across languages; a later import can promote one to the core list.
INSERT INTO concepts (gloss_en, core) VALUES ($1, $2)
ON CONFLICT (gloss_en, definition) DO UPDATE SET core = concepts.core OR excluded.core
RETURNING id;

-- name: CreateItem :one
INSERT INTO items (kind, language_code, status, author_id) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: FindEntry :one
SELECT id FROM entries
WHERE language_code = $1 AND headword = $2 AND part_of_speech IS NOT DISTINCT FROM $3;

-- name: CreateEntry :exec
INSERT INTO entries (id, language_code, headword, part_of_speech) VALUES ($1, $2, $3, $4);

-- name: FindSense :one
SELECT id FROM senses WHERE entry_id = $1 AND concept_id = $2;

-- name: CreateSense :exec
INSERT INTO senses (id, entry_id, concept_id, gloss, position)
VALUES ($1, $2, $3, $4, (SELECT coalesce(max(position), 0) + 1 FROM senses WHERE entry_id = $2));

-- name: LinkSource :exec
INSERT INTO item_sources (item_id, source_id, reference) VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;
