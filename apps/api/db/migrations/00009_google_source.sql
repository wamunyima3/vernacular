-- The register row Google Translate candidates cite (P1-AI-02, AI-1, AI-2).

-- +goose Up
INSERT INTO sources (kind, name, rights_holder, licence, url, allowed_uses, notes) VALUES (
    'machine', 'Google Translate', 'Google LLC', 'Google Cloud Terms of Service',
    'https://cloud.google.com/terms',
    '{redistribute,derive,commercial}',
    'Cloud Translation API (v2, NMT): one candidate per concept, bem and ny only. Stored as machine until reviewed; never used as training data.'
);

-- +goose Down
DELETE FROM items WHERE id IN (
    SELECT item_id FROM item_sources JOIN sources s ON s.id = source_id WHERE s.name = 'Google Translate'
) AND status = 'machine';
DELETE FROM item_sources WHERE source_id = (SELECT id FROM sources WHERE name = 'Google Translate');
DELETE FROM sources WHERE name = 'Google Translate';
