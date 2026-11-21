-- servable() is the one definition of which statuses the API may return at
-- all (ADR-0004): never `disputed`, never the internal `rejected`. Read
-- endpoints narrow it further with `min_status`.

-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION servable(status trust_status) RETURNS boolean
LANGUAGE sql IMMUTABLE AS $$ SELECT status NOT IN ('disputed', 'rejected') $$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION servable;
