-- Reviewer qualification (DQ-3, US-24): a would-be reviewer answers a short
-- test drawn from verified words, and a senior reviewer of the language
-- confirms a pass, which approves them. The expected answers stay here, never
-- in the test the user is sent.

-- +goose Up
CREATE TABLE qualification_attempts (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       text        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    language_code text        NOT NULL REFERENCES languages (code),
    expected      uuid[]      NOT NULL, -- the right entry for each question, in order
    score         int,
    passed        boolean,
    answered_at   timestamptz,
    confirmed_by  text        REFERENCES users (id),
    confirmed_at  timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    CHECK ((answered_at IS NULL) = (score IS NULL) AND (score IS NULL) = (passed IS NULL)),
    CHECK (confirmed_at IS NULL OR passed)
);
CREATE INDEX qualification_attempts_user_idx ON qualification_attempts (user_id, language_code, created_at);
CREATE INDEX qualification_attempts_pending_idx ON qualification_attempts (language_code, answered_at)
    WHERE passed AND confirmed_at IS NULL;
SELECT track('qualification_attempts');

-- +goose Down
DROP TABLE qualification_attempts;
