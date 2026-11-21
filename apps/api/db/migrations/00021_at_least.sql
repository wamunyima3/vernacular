-- at_least() keeps items at or above the caller's min_status (ADR-0004,
-- TRN-4, AI-6). It relies on the trust_status enum's declared order:
-- verified, reviewed, sourced, machine, then disputed and rejected (00014),
-- so with any minimum from the API those two are never servable. A new
-- status must be placed in that order on purpose.

-- +goose Up
CREATE FUNCTION at_least(status trust_status, min trust_status) RETURNS boolean
LANGUAGE sql IMMUTABLE RETURN status <= min;

-- +goose Down
DROP FUNCTION at_least;
