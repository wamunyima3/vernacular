-- name: UserExists :one
SELECT EXISTS (SELECT 1 FROM users WHERE id = $1);

-- name: RecordConsent :one
-- A new signed form; recording it again after a withdrawal starts consent anew.
INSERT INTO speakers (user_id, consent_key, consented_at, recorded_by, languages)
VALUES (@user_id, @consent_key, now(), @recorded_by, @languages::text[])
ON CONFLICT (user_id) DO UPDATE SET
    consent_key = EXCLUDED.consent_key, consented_at = now(), languages = EXCLUDED.languages,
    recorded_by = EXCLUDED.recorded_by, withdrawn_at = NULL, updated_at = now()
RETURNING *;

-- name: GetSpeaker :one
SELECT * FROM speakers WHERE user_id = $1;

-- name: WithdrawConsent :one
UPDATE speakers SET withdrawn_at = now(), updated_at = now()
WHERE user_id = $1 AND withdrawn_at IS NULL
RETURNING *;

-- name: EntryLanguage :one
SELECT language_code FROM entries WHERE id = $1;

-- name: SaveClip :one
INSERT INTO audio_clips (entry_id, speaker_id, content_type, size)
VALUES (@entry_id, @speaker_id, @content_type, @size)
ON CONFLICT (entry_id, speaker_id) DO UPDATE SET
    content_type = EXCLUDED.content_type, size = EXCLUDED.size, recorded_at = now()
RETURNING *;

-- name: ListClips :many
-- An entry's clips from speakers whose consent is in force, oldest first.
SELECT c.* FROM audio_clips c JOIN speakers s ON s.user_id = c.speaker_id
WHERE c.entry_id = $1 AND s.withdrawn_at IS NULL
ORDER BY c.recorded_at;

-- name: GetClip :one
SELECT c.* FROM audio_clips c JOIN speakers s ON s.user_id = c.speaker_id
WHERE c.id = $1 AND s.withdrawn_at IS NULL;

-- name: DeleteClip :exec
DELETE FROM audio_clips WHERE id = $1;

-- name: DeleteSpeakerClips :exec
DELETE FROM audio_clips WHERE speaker_id = $1;
