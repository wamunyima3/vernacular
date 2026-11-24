-- Dataset releases (PRV-4, US-29).

-- Newest first by number, so 2027.10 follows 2027.9.
-- name: ListReleases :many
SELECT * FROM releases ORDER BY split_part(name, '.', 1)::int DESC, split_part(name, '.', 2)::int DESC;

-- Serialises builds, so two admins cannot take the same name; reads go on.
-- name: LockReleases :exec
LOCK TABLE releases IN SHARE ROW EXCLUSIVE MODE;

-- <year>.<n>: the year's next number.
-- name: CreateRelease :one
INSERT INTO releases (name, published_at, changelog)
SELECT to_char(now(), 'YYYY') || '.' || (count(*) + 1), now(), @changelog
FROM releases WHERE name LIKE to_char(now(), 'YYYY') || '.%'
RETURNING *;

-- Copies every verified item as it is now. A snapshot is the item's own row
-- without generated columns or timestamps, so it changes only with content.
-- name: FreezeVerified :execrows
INSERT INTO release_items (release_id, item_id, snapshot)
SELECT @release_id::uuid, i.id,
       jsonb_build_object('kind', i.kind, 'language', i.language_code) || coalesce(CASE i.kind
           WHEN 'entry' THEN (SELECT to_jsonb(e) - '{id,headword_fold,created_at,updated_at}'::text[]
                              FROM entries e WHERE e.id = i.id)
           WHEN 'sense' THEN (SELECT to_jsonb(s) - '{id,created_at,updated_at}'::text[]
                              FROM senses s WHERE s.id = i.id)
           WHEN 'example' THEN (SELECT to_jsonb(x) - '{id,created_at,updated_at}'::text[]
                                FROM examples x WHERE x.id = i.id)
           WHEN 'phrase' THEN (SELECT to_jsonb(p) - '{id,text_fold,source_en_fold,created_at,updated_at}'::text[]
                               FROM phrases p WHERE p.id = i.id)
       END, '{}')
FROM items i
WHERE i.status = 'verified';

-- Counts per language and kind, and the changes against the previous release
-- (none given: everything is added).
-- name: SummariseRelease :one
WITH cur AS (SELECT item_id, snapshot FROM release_items WHERE release_id = @release_id::uuid),
     prev AS (SELECT item_id, snapshot FROM release_items WHERE release_id = sqlc.narg(previous_id)::uuid),
     per_kind AS (
         SELECT snapshot ->> 'language' AS lang, snapshot ->> 'kind' AS kind, count(*) AS n
         FROM cur GROUP BY 1, 2
     ),
     diff AS (
         SELECT count(*) FILTER (WHERE prev.item_id IS NULL) AS added,
                count(*) FILTER (WHERE cur.snapshot <> prev.snapshot) AS changed,
                count(*) FILTER (WHERE cur.item_id IS NULL) AS removed
         FROM cur FULL JOIN prev ON prev.item_id = cur.item_id
     )
UPDATE releases SET counts = jsonb_build_object(
    'languages', (SELECT coalesce(jsonb_object_agg(lang, kinds), '{}')
                  FROM (SELECT lang, jsonb_object_agg(kind, n) AS kinds FROM per_kind GROUP BY lang) l),
    'added', diff.added, 'changed', diff.changed, 'removed', diff.removed
)
FROM diff
WHERE id = @release_id::uuid
RETURNING releases.*;

-- name: MaterialiseRelease :exec
SELECT materialise_release(@release_id::uuid, @schema::text);

-- The release a pinned read names; `latest` is the newest.
-- name: ResolveRelease :one
SELECT name FROM releases WHERE @name::text IN ('latest', name)
ORDER BY split_part(name, '.', 1)::int DESC, split_part(name, '.', 2)::int DESC
LIMIT 1;
