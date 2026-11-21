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
