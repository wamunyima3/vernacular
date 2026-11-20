-- Source register rules (P1-SRC-01, DQ-1).

-- +goose Up
ALTER TABLE sources ADD CONSTRAINT sources_allowed_uses_check
    CHECK (allowed_uses <@ ARRAY['redistribute', 'derive', 'commercial', 'train']); -- contract AllowedUse

-- +goose Down
ALTER TABLE sources DROP CONSTRAINT sources_allowed_uses_check;
