-- The register row community contributions cite (P1-REV-01, CON-1, DQ-1).
-- The outward licence waits for the dataset licence decision (open question
-- in docs/requirements.md), so no uses are granted yet.

-- +goose Up
INSERT INTO sources (kind, name, rights_holder, licence, allowed_uses, notes) VALUES (
    'contributor', 'Community contributors', 'Each contributor',
    'Contributor agreement (draft): published under the dataset licence',
    '{}',
    'Set licence and allowed uses once the dataset licence is decided. Each item records its author.'
);

-- +goose Down
DELETE FROM items WHERE id IN (
    SELECT item_id FROM item_sources JOIN sources s ON s.id = source_id WHERE s.name = 'Community contributors'
) AND author_id IS NOT NULL;
DELETE FROM item_sources WHERE source_id = (SELECT id FROM sources WHERE name = 'Community contributors');
DELETE FROM sources WHERE name = 'Community contributors';
