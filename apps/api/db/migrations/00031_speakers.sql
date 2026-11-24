-- Speakers who agreed in writing to have their voice recorded (LEX-8, NFR-6,
-- US-33). Voice is personal data under the Data Protection Act, 2021, so a
-- clip may only come from a speaker with a signed consent form on file and
-- not withdrawn. The form is kept in the object store like permission
-- documents; withdrawing keeps this row as the record of both events.

-- +goose Up
CREATE TABLE speakers (
    user_id      text        PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    consent_key  text        NOT NULL, -- the signed form, in the object store
    consented_at timestamptz NOT NULL,
    recorded_by  text        NOT NULL REFERENCES users (id),
    withdrawn_at timestamptz,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
SELECT track('speakers');

-- +goose Down
DROP TABLE speakers;
