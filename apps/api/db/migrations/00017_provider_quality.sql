-- Acceptance and correction rates per AI provider and language (AI-5,
-- P1-REV-06): review decisions on the items each machine source proposed.
-- An item two providers proposed counts for both. Reviews a send-back set
-- aside do not count, as for promotion.

-- +goose Up
CREATE VIEW provider_quality AS
SELECT s.name AS provider, i.language_code,
       count(*)::int AS decided,
       count(*) FILTER (WHERE r.decision = 'accept')::int AS accepted,
       count(*) FILTER (WHERE r.decision = 'correct')::int AS corrected
FROM reviews r
JOIN items i ON i.id = r.item_id
JOIN item_sources isrc ON isrc.item_id = i.id
JOIN sources s ON s.id = isrc.source_id AND s.kind = 'machine'
WHERE r.set_aside_by IS NULL
GROUP BY s.name, i.language_code;

-- +goose Down
DROP VIEW provider_quality;
