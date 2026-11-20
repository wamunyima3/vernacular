-- ADR-0004 amendment: items two reviewers reject leave the queue as
-- `rejected`. Never returned by the API, so the contract is unchanged.

-- +goose Up
ALTER TYPE trust_status ADD VALUE 'rejected';

-- +goose Down
-- Postgres cannot drop an enum value; move any rejected items back to the
-- queue instead, which is what the earlier rules did.
UPDATE items i SET status = CASE
    WHEN EXISTS (SELECT 1 FROM item_sources isrc JOIN sources s ON s.id = isrc.source_id
                 WHERE isrc.item_id = i.id AND s.kind <> 'machine') THEN 'sourced'::trust_status
    ELSE 'machine'::trust_status END
WHERE status = 'rejected';
