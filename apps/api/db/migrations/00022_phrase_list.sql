-- The English source list of the phrase packs (PHR-1, P1-API-07b): what
-- each pack should say. Each line points to the concept it expresses, so
-- candidates and reviews can link a language's phrase to it.

-- +goose Up
CREATE TABLE phrase_list (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    pack       text        NOT NULL CHECK (pack ~ '^[a-z][a-z0-9-]*$'),
    source_en  text        NOT NULL CHECK (source_en <> ''),
    notes      text        NOT NULL DEFAULT '',
    concept_id uuid        NOT NULL REFERENCES concepts (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (pack, source_en)
);
SELECT track('phrase_list');

-- +goose Down
DROP TABLE phrase_list;
