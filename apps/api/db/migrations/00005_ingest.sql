-- Ingestion support (P0-CNT-03a, P1-SRC-02).

-- +goose Up

-- Sources such as GATITOS give no part of speech. Reviewers add it.
ALTER TABLE entries ALTER COLUMN part_of_speech DROP NOT NULL;

-- Keys that make re-running an import idempotent.
ALTER TABLE entries ADD CONSTRAINT entries_identity_key
    UNIQUE NULLS NOT DISTINCT (language_code, headword, part_of_speech);
ALTER TABLE senses ADD CONSTRAINT senses_entry_concept_key UNIQUE (entry_id, concept_id);

-- GATITOS: Google's word and short-phrase lists, part of SMOL (docs/sources/smol-gatitos.md).
INSERT INTO sources (kind, name, rights_holder, licence, url, allowed_uses, notes) VALUES (
    'licensed', 'GATITOS', 'Google LLC', 'CC BY 4.0',
    'https://huggingface.co/datasets/google/smol',
    '{redistribute,derive,commercial}',
    'Attribution: Caswell et al., SMOL, 2025. Bemba only among our languages.'
);

-- +goose Down
-- Imported rows go with their source.
DELETE FROM items WHERE id IN (
    SELECT item_id FROM item_sources WHERE source_id = (SELECT id FROM sources WHERE name = 'GATITOS')
);
DELETE FROM sources WHERE name = 'GATITOS';
ALTER TABLE senses DROP CONSTRAINT senses_entry_concept_key;
ALTER TABLE entries DROP CONSTRAINT entries_identity_key;
-- Rows imported without a part of speech would block NOT NULL, so they go first.
DELETE FROM items WHERE id IN (SELECT id FROM entries WHERE part_of_speech IS NULL);
ALTER TABLE entries ALTER COLUMN part_of_speech SET NOT NULL;
