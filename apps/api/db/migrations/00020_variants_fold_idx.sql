-- Exact variant lookup in entry search (P1-API-02); the trigram index only
-- serves similarity.

-- +goose Up
CREATE INDEX variants_fold_idx ON variants (fold);

-- +goose Down
DROP INDEX variants_fold_idx;
