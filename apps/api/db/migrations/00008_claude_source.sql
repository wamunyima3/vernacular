-- The register row that AI candidates cite (P1-AI-01, AI-1, AI-2). Each
-- candidate's item_sources row records the exact model and the date.

-- +goose Up
INSERT INTO sources (kind, name, rights_holder, licence, url, allowed_uses, notes) VALUES (
    'machine', 'Claude', 'Anthropic PBC', 'Anthropic Commercial Terms of Service',
    'https://www.anthropic.com/legal/commercial-terms',
    '{redistribute,derive,commercial}',
    'Candidates generated per concept. Stored as machine until reviewed; never used as training data.'
);

-- +goose Down
DELETE FROM items WHERE id IN (
    SELECT item_id FROM item_sources JOIN sources s ON s.id = source_id WHERE s.name = 'Claude'
) AND status = 'machine';
DELETE FROM item_sources WHERE source_id = (SELECT id FROM sources WHERE name = 'Claude');
DELETE FROM sources WHERE name = 'Claude';
