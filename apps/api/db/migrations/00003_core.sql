-- Core lexicon, review and release schema (P0-ENG-01, ADR-0006).
--
-- Every reviewable row (entry, sense, example, phrase) shares one row in
-- `items`, which holds its trust status. Reviews, source links, release
-- snapshots and the audit log point at `items`, so they need one foreign key
-- instead of one per table.
--
-- User IDs are text: the auth task (P0-ENG-06) adds the users table and the
-- foreign keys to it.

-- +goose Up

-- Shared helpers --------------------------------------------------------------

-- fold() is the single definition of search folding (NFR-7): NFC, lower case,
-- curly apostrophes made straight, ŋ spelled ng' (style guides), accents
-- removed. Go's textnorm.Fold must give the same result. unaccent() is only
-- STABLE, so it is called with an explicit dictionary to allow use in indexes
-- and generated columns.
-- +goose StatementBegin
CREATE FUNCTION fold(t text) RETURNS text
LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT
RETURN public.unaccent(
    'public.unaccent'::regdictionary,
    replace(translate(lower(normalize(t, NFC)), '’ʼ‘', ''''''''), 'ŋ', 'ng''')
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END $$;
-- +goose StatementEnd

-- Every write to an audited table lands here (NFR-8). The acting user comes
-- from `SET LOCAL app.actor_id = '...'` inside the writing transaction.
CREATE TABLE audit_log (
    id         bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    table_name text        NOT NULL,
    row_id     text        NOT NULL,
    action     text        NOT NULL CHECK (action IN ('insert', 'update', 'delete')),
    actor_id   text,
    before     jsonb,
    after      jsonb,
    at         timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_log_row_idx ON audit_log (table_name, row_id, at);

-- The trigger arguments name the key columns; row_id joins their values with ':'.
-- +goose StatementBegin
CREATE FUNCTION audit() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    old_row jsonb := CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE to_jsonb(OLD) END;
    new_row jsonb := CASE WHEN TG_OP = 'DELETE' THEN NULL ELSE to_jsonb(NEW) END;
BEGIN
    INSERT INTO audit_log (table_name, row_id, action, actor_id, before, after)
    VALUES (
        TG_TABLE_NAME,
        array_to_string(ARRAY(SELECT coalesce(new_row, old_row) ->> k FROM unnest(TG_ARGV) AS k), ':'),
        lower(TG_OP),
        nullif(current_setting('app.actor_id', true), ''),
        old_row,
        new_row
    );
    RETURN NULL;
END $$;
-- +goose StatementEnd

-- Applies the updated_at and audit triggers to a table. Every table created
-- below calls it, so no table can be forgotten. `key` lists the primary key
-- columns, so composite keys need no special column.
-- +goose StatementBegin
CREATE FUNCTION track(tbl regclass, key text[] DEFAULT '{id}') RETURNS void LANGUAGE plpgsql AS $$
BEGIN
    EXECUTE format('CREATE TRIGGER set_updated_at BEFORE UPDATE ON %s FOR EACH ROW EXECUTE FUNCTION set_updated_at()', tbl);
    EXECUTE format('CREATE TRIGGER audit AFTER INSERT OR UPDATE OR DELETE ON %s FOR EACH ROW EXECUTE FUNCTION audit(%s)',
        tbl, (SELECT string_agg(quote_literal(k), ', ') FROM unnest(key) AS k));
END $$;
-- +goose StatementEnd

SELECT track('languages', '{code}');

-- Trust (ADR-0004) -------------------------------------------------------------

CREATE TYPE trust_status AS ENUM ('verified', 'reviewed', 'sourced', 'machine', 'disputed');
CREATE TYPE item_kind AS ENUM ('entry', 'sense', 'example', 'phrase');

CREATE TABLE items (
    id            uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    kind          item_kind    NOT NULL,
    language_code text         NOT NULL REFERENCES languages (code),
    status        trust_status NOT NULL,
    author_id     text,
    created_at    timestamptz  NOT NULL DEFAULT now(),
    updated_at    timestamptz  NOT NULL DEFAULT now()
);
CREATE INDEX items_queue_idx ON items (language_code, status, kind);
SELECT track('items');

-- Sources (DQ-1) ----------------------------------------------------------------

CREATE TYPE source_kind AS ENUM ('licensed', 'contributor', 'machine');

CREATE TABLE sources (
    id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    kind           source_kind NOT NULL,
    name           text        NOT NULL UNIQUE,
    rights_holder  text        NOT NULL,
    licence        text        NOT NULL,
    url            text,
    permission_key text,                        -- object storage key of the signed permission document
    allowed_uses   text[]      NOT NULL DEFAULT '{}',
    notes          text        NOT NULL DEFAULT '',
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);
SELECT track('sources');

CREATE TABLE item_sources (
    item_id    uuid        NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    source_id  uuid        NOT NULL REFERENCES sources (id),
    reference  text        NOT NULL DEFAULT '',  -- page, line or row in the source
    model      text,                             -- model version, for machine sources (AI-2)
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (item_id, source_id, reference)
);

-- Lexicon (LEX-1…7) ---------------------------------------------------------------

CREATE TABLE varieties (
    language_code text NOT NULL REFERENCES languages (code),
    code          text NOT NULL CHECK (code ~ '^[a-z][a-z0-9-]*$'),
    name          text NOT NULL,
    PRIMARY KEY (language_code, code)
);

INSERT INTO varieties (language_code, code, name) VALUES
    ('bem', 'copperbelt', 'Copperbelt urban'),
    ('bem', 'luapula',    'Luapula'),
    ('bem', 'northern',   'Northern'),
    ('nya', 'town',       'Lusaka town Nyanja'),
    ('nya', 'chewa',      'Eastern Province Chewa');

CREATE TABLE concepts (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    gloss_en   text        NOT NULL,
    definition text        NOT NULL DEFAULT '',
    domain     text        NOT NULL DEFAULT 'general',
    core       boolean     NOT NULL DEFAULT false,  -- part of the ~4,000 core list (LEX-7)
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (gloss_en, definition)
);
CREATE INDEX concepts_gloss_trgm_idx ON concepts USING gin (fold(gloss_en) gin_trgm_ops);
SELECT track('concepts');

CREATE TYPE part_of_speech AS ENUM (
    'noun', 'verb', 'adjective', 'adverb', 'pronoun', 'numeral', 'conjunction',
    'preposition', 'interjection', 'phrase', 'ideophone'
);
CREATE TYPE register AS ENUM (
    'neutral', 'formal', 'informal', 'town', 'respectful', 'vulgar', 'archaic'
);

CREATE TABLE entries (
    id             uuid           PRIMARY KEY REFERENCES items (id) ON DELETE CASCADE,
    language_code  text           NOT NULL REFERENCES languages (code),
    headword       text           NOT NULL CHECK (headword = normalize(headword, NFC)),
    headword_fold  text           GENERATED ALWAYS AS (fold(headword)) STORED,
    part_of_speech part_of_speech NOT NULL,
    noun_class     text           CHECK (noun_class ~ '^[0-9]{1,2}(/[0-9]{1,2})?$'),
    plural         text,
    root           text,
    variety        text,
    register       register       NOT NULL DEFAULT 'neutral',
    created_at     timestamptz    NOT NULL DEFAULT now(),
    updated_at     timestamptz    NOT NULL DEFAULT now(),
    FOREIGN KEY (language_code, variety) REFERENCES varieties (language_code, code)
);
CREATE INDEX entries_headword_trgm_idx ON entries USING gin (headword_fold gin_trgm_ops);
CREATE INDEX entries_lookup_idx ON entries (language_code, headword_fold);
CREATE INDEX entries_root_idx ON entries (language_code, root);
SELECT track('entries');

CREATE TABLE variants (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id   uuid        NOT NULL REFERENCES entries (id) ON DELETE CASCADE,
    spelling   text        NOT NULL CHECK (spelling = normalize(spelling, NFC)),
    fold       text        GENERATED ALWAYS AS (fold(spelling)) STORED,
    note       text        NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (entry_id, spelling)
);
CREATE INDEX variants_fold_trgm_idx ON variants USING gin (fold gin_trgm_ops);
SELECT track('variants');

CREATE TABLE senses (
    id         uuid        PRIMARY KEY REFERENCES items (id) ON DELETE CASCADE,
    entry_id   uuid        NOT NULL REFERENCES entries (id) ON DELETE CASCADE,
    concept_id uuid        REFERENCES concepts (id),
    position   smallint    NOT NULL DEFAULT 1,
    gloss      text        NOT NULL,
    definition text        NOT NULL DEFAULT '',
    notes      text        NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX senses_entry_idx ON senses (entry_id, position);
CREATE INDEX senses_concept_idx ON senses (concept_id);
SELECT track('senses');

CREATE TABLE examples (
    id             uuid        PRIMARY KEY REFERENCES items (id) ON DELETE CASCADE,
    sense_id       uuid        NOT NULL REFERENCES senses (id) ON DELETE CASCADE,
    text           text        NOT NULL,
    translation_en text        NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX examples_sense_idx ON examples (sense_id);
SELECT track('examples');

-- Phrases (PHR-1…3) -----------------------------------------------------------------

CREATE TABLE phrases (
    id             uuid        PRIMARY KEY REFERENCES items (id) ON DELETE CASCADE,
    language_code  text        NOT NULL REFERENCES languages (code),
    concept_id     uuid        REFERENCES concepts (id),
    pack           text        NOT NULL CHECK (pack ~ '^[a-z][a-z0-9-]*$'),
    text           text        NOT NULL CHECK (text = normalize(text, NFC)),
    text_fold      text        GENERATED ALWAYS AS (fold(text)) STORED,
    source_en      text        NOT NULL,
    source_en_fold text        GENERATED ALWAYS AS (fold(source_en)) STORED,
    notes          text        NOT NULL DEFAULT '',
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX phrases_pack_idx ON phrases (language_code, pack);
CREATE INDEX phrases_exact_en_idx ON phrases (language_code, source_en_fold);
CREATE INDEX phrases_exact_text_idx ON phrases (language_code, text_fold);
SELECT track('phrases');

-- Review (CON-2…4, DQ-2) -----------------------------------------------------------

CREATE TYPE review_decision AS ENUM ('accept', 'correct', 'reject');

CREATE TABLE reviews (
    id          uuid            PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id     uuid            NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    reviewer_id text            NOT NULL,
    decision    review_decision NOT NULL,
    correction  jsonb,                      -- the changed fields, for decision = 'correct'
    comment     text            NOT NULL DEFAULT '',
    created_at  timestamptz     NOT NULL DEFAULT now(),
    updated_at  timestamptz     NOT NULL DEFAULT now(),
    UNIQUE (item_id, reviewer_id),          -- one decision per reviewer per item
    CHECK ((decision = 'correct') = (correction IS NOT NULL))
);
CREATE INDEX reviews_reviewer_idx ON reviews (reviewer_id, created_at);
SELECT track('reviews');

-- Releases (PRV-4) ------------------------------------------------------------------

CREATE TABLE releases (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    name         text        NOT NULL UNIQUE CHECK (name ~ '^[0-9]{4}\.[0-9]+$'),
    published_at timestamptz,
    counts       jsonb       NOT NULL DEFAULT '{}',
    changelog    text        NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
SELECT track('releases');

-- A frozen copy of each verified item as it was when the release was built.
CREATE TABLE release_items (
    release_id uuid  NOT NULL REFERENCES releases (id) ON DELETE CASCADE,
    item_id    uuid  NOT NULL REFERENCES items (id),
    snapshot   jsonb NOT NULL,
    PRIMARY KEY (release_id, item_id)
);

-- +goose Down
DROP TABLE release_items;
DROP TABLE releases;
DROP TABLE reviews;
DROP TYPE review_decision;
DROP TABLE phrases;
DROP TABLE examples;
DROP TABLE senses;
DROP TABLE variants;
DROP TABLE entries;
DROP TYPE register;
DROP TYPE part_of_speech;
DROP TABLE concepts;
DROP TABLE varieties;
DROP TABLE item_sources;
DROP TABLE sources;
DROP TYPE source_kind;
DROP TABLE items;
DROP TYPE item_kind;
DROP TYPE trust_status;
DROP TRIGGER audit ON languages;
DROP TRIGGER set_updated_at ON languages;
DROP FUNCTION track(regclass, text[]);
DROP FUNCTION audit();
DROP TABLE audit_log;
DROP FUNCTION set_updated_at();
DROP FUNCTION fold(text);
