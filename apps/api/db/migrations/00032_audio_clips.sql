-- Pronunciation clips (LEX-8, US-33). A speaker records entries in the
-- languages their consent form names, one clip per entry; recording again
-- replaces it. The file lives in the object store at audio/{entry}/{speaker}.
-- Withdrawing consent deletes the speaker's clips and files (NFR-6).

-- +goose Up
ALTER TABLE speakers ADD COLUMN languages text[] NOT NULL DEFAULT '{}';

CREATE TABLE audio_clips (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id     uuid        NOT NULL REFERENCES entries (id) ON DELETE CASCADE,
    speaker_id   text        NOT NULL REFERENCES speakers (user_id) ON DELETE CASCADE,
    content_type text        NOT NULL,
    size         int         NOT NULL CHECK (size > 0),
    recorded_at  timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    UNIQUE (entry_id, speaker_id)
);
SELECT track('audio_clips');

-- +goose Down
DROP TABLE audio_clips;
ALTER TABLE speakers DROP COLUMN languages;
