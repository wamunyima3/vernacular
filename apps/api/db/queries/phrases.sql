-- name: ListPack :many
-- A phrase pack in one language (PHR-1, US-08): its phrases at or above
-- min_status, in English order.
SELECT p.id, p.text, p.source_en, p.notes, p.concept_id
FROM phrases p JOIN items i ON i.id = p.id
WHERE p.language_code = @language AND p.pack = @pack AND at_least(i.status, @min_status)
ORDER BY p.source_en_fold, p.id
LIMIT 1000;

-- name: UpsertListLine :one
-- Adds or updates a line of the English source list; unchanged lines are not
-- written, so re-runs add no audit rows. Reports whether the line is new.
INSERT INTO phrase_list (pack, source_en, notes, concept_id)
VALUES (@pack, @source_en, @notes, @concept_id)
ON CONFLICT (pack, source_en) DO UPDATE SET notes = excluded.notes, updated_at = now()
WHERE phrase_list.notes <> excluded.notes
RETURNING (xmax = 0) AS inserted;

-- name: MatchEnglish :many
-- Verified phrases whose English source equals the text once folded (TRN-2).
SELECT p.id, p.text, p.source_en, p.notes, p.concept_id
FROM phrases p JOIN items i ON i.id = p.id
WHERE p.language_code = @language AND p.source_en_fold = fold(@text::text) AND i.status = 'verified'
ORDER BY p.pack, p.id
LIMIT 10;

-- name: MatchText :many
-- Verified phrases whose own text equals the text once folded (TRN-2).
SELECT p.id, p.text, p.source_en, p.notes, p.concept_id
FROM phrases p JOIN items i ON i.id = p.id
WHERE p.language_code = @language AND p.text_fold = fold(@text::text) AND i.status = 'verified'
ORDER BY p.pack, p.id
LIMIT 10;

-- name: ListLineIDs :many
SELECT id FROM phrase_list ORDER BY pack, source_en;

-- name: GetListLine :one
SELECT id, pack, source_en, notes, concept_id FROM phrase_list WHERE id = $1;

-- Serialises candidate writes per language, so two providers proposing the
-- same phrase at once agree on one item.
-- name: LockPhrases :exec
SELECT pg_advisory_xact_lock(hashtext('phrases:' || @language_code::text));

-- name: FindPhrase :one
-- The phrase already proposed for this line with this text, folded.
SELECT id FROM phrases
WHERE language_code = @language AND pack = @pack AND source_en_fold = fold(@source_en::text) AND text_fold = fold(@text::text)
LIMIT 1;

-- name: InsertPhrase :exec
INSERT INTO phrases (id, language_code, concept_id, pack, text, source_en, notes)
VALUES (@id, @language, @concept_id, @pack, @text, @source_en, @notes);

-- name: ListPhrasesByID :many
SELECT p.id, p.text, p.source_en, p.notes, p.concept_id
FROM phrases p WHERE p.id = ANY(@ids::uuid[]);

-- name: GetPhraseLine :one
-- The source-list line a phrase translates.
SELECT l.id, l.pack, l.source_en, l.notes, l.concept_id
FROM phrases p JOIN phrase_list l ON l.pack = p.pack AND l.source_en = p.source_en
WHERE p.id = $1;
