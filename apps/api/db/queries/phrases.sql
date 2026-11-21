-- name: ListPack :many
-- A phrase pack in one language (PHR-1, US-08): its phrases at or above
-- min_status, in English order.
SELECT p.id, p.text, p.source_en, p.notes, p.concept_id
FROM phrases p JOIN items i ON i.id = p.id
WHERE p.language_code = @language AND p.pack = @pack AND at_least(i.status, @min_status)
ORDER BY p.source_en_fold, p.id
LIMIT 1000;
