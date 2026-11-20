-- Search stems (P1-SRC-04, LEX-4). Computed in Go (lexicon.Stem) on write;
-- `ingest -restem` recomputes them when the rules change.

-- +goose Up
ALTER TABLE entries ADD COLUMN stem text;
CREATE INDEX entries_stem_idx ON entries (language_code, stem);

-- +goose Down
ALTER TABLE entries DROP COLUMN stem;
