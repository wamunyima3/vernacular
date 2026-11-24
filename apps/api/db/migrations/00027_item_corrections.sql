-- A contribution can correct an existing item (CON-1, US-18): the proposal
-- is a new candidate, linked to the item it would replace, so the original
-- keeps its history and stays served until reviewers verify the correction.

-- +goose Up
ALTER TABLE items ADD COLUMN corrects uuid REFERENCES items (id);
CREATE INDEX items_corrects_idx ON items (corrects) WHERE corrects IS NOT NULL;

-- +goose Down
ALTER TABLE items DROP COLUMN corrects;
