-- name: ListTrainingText :many
-- What the language detector learns from (ADR-0010): each active language's
-- servable headwords, examples and phrases, and English from concept glosses
-- and example translations.
SELECT e.language_code AS language, e.headword AS text
FROM entries e JOIN items i ON i.id = e.id JOIN languages l ON l.code = e.language_code
WHERE l.status = 'active' AND servable(i.status)
UNION ALL
SELECT i.language_code, x.text
FROM examples x JOIN items i ON i.id = x.id JOIN languages l ON l.code = i.language_code
WHERE l.status = 'active' AND servable(i.status)
UNION ALL
SELECT p.language_code, p.text
FROM phrases p JOIN items i ON i.id = p.id JOIN languages l ON l.code = p.language_code
WHERE l.status = 'active' AND servable(i.status)
UNION ALL
SELECT 'eng', gloss_en FROM concepts
UNION ALL
SELECT 'eng', x.translation_en FROM examples x;
