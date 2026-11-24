-- retired() is true for an item a verified correction replaces (US-18,
-- P1-REV-01b2). The item keeps its row, status and history, but read
-- endpoints and the review queue leave it out. It is derived, so undoing the
-- accept that verified a correction brings the original back.
-- The body is a string so `items` resolves at call time: in a pinned
-- release's schema it reads that release's items.

-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION retired(item uuid) RETURNS boolean
LANGUAGE sql STABLE AS $$
    SELECT EXISTS (SELECT 1 FROM items c WHERE c.corrects = item AND c.status = 'verified')
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION retired;
