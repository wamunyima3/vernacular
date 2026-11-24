-- Re-check batches (DQ-6, P2-10b). An item's batch is its first source and
-- model: the import or the AI run it came from. When more than 2% of a
-- batch's re-checks in a quarter fail, an admin sends the batch's verified
-- items back to review, recorded here once per quarter.

-- +goose Up
-- One row per item: its first source and model. Filtering on item_id is
-- pushed into the DISTINCT ON, so a join reads only the items it needs.
CREATE VIEW item_batches AS
SELECT DISTINCT ON (item_id) item_id, source_id, coalesce(model, '') AS model
FROM item_sources
ORDER BY item_id, created_at, source_id;

CREATE TABLE recheck_batches (
    quarter       text        NOT NULL,
    language_code text        NOT NULL REFERENCES languages (code),
    source_id     uuid        NOT NULL REFERENCES sources (id),
    model         text        NOT NULL DEFAULT '',
    sent_back_by  text        NOT NULL,
    items         int         NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (quarter, language_code, source_id, model)
);
SELECT track('recheck_batches');

-- +goose Down
DROP TABLE recheck_batches;
DROP VIEW item_batches;
