-- name: UserExists :one
SELECT EXISTS (SELECT 1 FROM users WHERE id = $1);

-- name: RecordConsent :one
-- A new signed form; recording it again after a withdrawal starts consent anew.
INSERT INTO speakers (user_id, consent_key, consented_at, recorded_by)
VALUES (@user_id, @consent_key, now(), @recorded_by)
ON CONFLICT (user_id) DO UPDATE SET
    consent_key = EXCLUDED.consent_key, consented_at = now(),
    recorded_by = EXCLUDED.recorded_by, withdrawn_at = NULL, updated_at = now()
RETURNING *;

-- name: GetSpeaker :one
SELECT * FROM speakers WHERE user_id = $1;

-- name: WithdrawConsent :one
UPDATE speakers SET withdrawn_at = now(), updated_at = now()
WHERE user_id = $1 AND withdrawn_at IS NULL
RETURNING *;
