-- has_verified_example() is true for a sense with at least one verified
-- example (LEX-5, P2-09c), for example coverage. The body is a string so
-- `examples` and `items` resolve at call time, as with retired().

-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION has_verified_example(sense uuid) RETURNS boolean
LANGUAGE sql STABLE AS $$
    SELECT EXISTS (SELECT 1 FROM examples x JOIN items i ON i.id = x.id
                   WHERE x.sense_id = sense AND i.status = 'verified')
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION has_verified_example;
