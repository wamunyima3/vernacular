-- name: GetSessionUser :one
-- The user behind a live session token, with their review approvals, in one
-- round trip. Banned users have no valid session.
SELECT u.id, u.name, u.email, u.role, (u.credit_consent_at IS NOT NULL)::boolean AS credited,
       coalesce(
           json_agg(json_build_object('language', a.language_code, 'level', a.level,
                                      'suspended', below_agreement(g.decided, g.agreed))
                    ORDER BY a.language_code)
               FILTER (WHERE a.user_id IS NOT NULL),
           '[]'
       )::jsonb AS approvals
FROM sessions s
JOIN users u ON u.id = s.user_id
LEFT JOIN reviewer_approvals a ON a.user_id = u.id
LEFT JOIN LATERAL (
    SELECT decided, agreed FROM reviewer_agreement
    WHERE user_id = a.user_id AND language_code = a.language_code
) g ON true
WHERE s.token = $1
  AND s.expires_at > now()
  AND NOT (coalesce(u.banned, false) AND (u.ban_expires IS NULL OR u.ban_expires > now()))
GROUP BY u.id;

-- name: SetRole :execrows
UPDATE users SET role = $2 WHERE email = $1;

-- name: ListReviewers :many
-- A language's reviewers with their agreement (DQ-7), suspended ones first.
SELECT u.id, u.name, a.level, g.decided, g.agreed,
       below_agreement(g.decided, g.agreed) AS suspended, a.retrained_at
FROM reviewer_approvals a
JOIN users u ON u.id = a.user_id
JOIN reviewer_agreement g ON g.user_id = a.user_id AND g.language_code = a.language_code
WHERE a.language_code = $1
ORDER BY suspended DESC, u.name, u.id;

-- name: MarkRetrained :execrows
UPDATE reviewer_approvals SET retrained_at = now() WHERE user_id = $1 AND language_code = $2;

-- name: ListApprovals :many
SELECT u.id, u.name, u.email, a.language_code, a.level,
       below_agreement(g.decided, g.agreed) AS suspended
FROM reviewer_approvals a
JOIN users u ON u.id = a.user_id
JOIN reviewer_agreement g ON g.user_id = a.user_id AND g.language_code = a.language_code
ORDER BY a.language_code, u.name, u.id;

-- name: UpsertApproval :one
-- Approves the user with the email, or changes their level. No row when no
-- user has the email.
INSERT INTO reviewer_approvals (user_id, language_code, level, approved_by)
SELECT u.id, @language_code, @level, @approved_by
FROM users u WHERE lower(u.email) = lower(@email)
ON CONFLICT (user_id, language_code) DO UPDATE SET level = EXCLUDED.level, approved_by = EXCLUDED.approved_by
RETURNING user_id;

-- name: DeleteApproval :execrows
DELETE FROM reviewer_approvals WHERE user_id = $1 AND language_code = $2;

-- name: SetCreditConsent :exec
-- Records when the user chose to be credited by name, or clears it (NFR-6).
UPDATE users SET credit_consent_at = CASE WHEN @credited::boolean THEN coalesce(credit_consent_at, now()) END
WHERE id = @id;

-- name: CreateKey :one
-- At most 10 live keys per user; no row means the user has that many.
INSERT INTO api_keys (user_id, name, prefix, hash)
SELECT @user_id, @name, @prefix, @hash
WHERE (SELECT count(*) FROM api_keys WHERE user_id = @user_id AND revoked_at IS NULL) < 10
RETURNING id, created_at;

-- name: ListKeys :many
-- With each key's calls today and in the last 30 days, UTC.
SELECT k.id, k.name, k.prefix, k.last_used_at, k.revoked_at, k.created_at,
    coalesce(sum(u.calls) FILTER (WHERE u.day = (now() AT TIME ZONE 'UTC')::date), 0)::int AS calls_today,
    coalesce(sum(u.calls), 0)::int AS calls_30d
FROM api_keys k
LEFT JOIN key_usage u ON u.key_id = k.id AND u.day > (now() AT TIME ZONE 'UTC')::date - 30
WHERE k.user_id = $1
GROUP BY k.id
ORDER BY k.revoked_at IS NOT NULL, k.created_at DESC
LIMIT 50;

-- name: RevokeKey :execrows
UPDATE api_keys SET revoked_at = now(), updated_at = now()
WHERE id = @id AND user_id = @user_id AND revoked_at IS NULL;

-- name: GetKeyOwner :one
-- The live key with this hash, and its owner unless banned.
SELECT k.id, k.user_id, k.last_used_at
FROM api_keys k JOIN users u ON u.id = k.user_id
WHERE k.hash = $1 AND k.revoked_at IS NULL AND NOT coalesce(u.banned, false);

-- name: TouchKey :exec
-- last_used_at is kept to the hour, so busy keys do not write on every call.
UPDATE api_keys SET last_used_at = now()
WHERE id = $1 AND (last_used_at IS NULL OR last_used_at < now() - interval '1 hour');

-- name: LatestQualification :one
-- The user's latest attempt in a language, to decide whether they may start another.
SELECT passed, confirmed_at, created_at FROM qualification_attempts
WHERE user_id = @user_id AND language_code = @language_code
ORDER BY created_at DESC LIMIT 1;

-- name: QualificationPool :many
-- Up to 60 verified, current meanings of a language at random, one per
-- entry, to draw questions and wrong options from (DQ-3).
SELECT entry_id, headword, gloss_en FROM (
    SELECT DISTINCT ON (e.id) e.id AS entry_id, e.headword, coalesce(c.gloss_en, s.gloss) AS gloss_en
    FROM senses s
    JOIN items i ON i.id = s.id
    JOIN entries e ON e.id = s.entry_id
    LEFT JOIN concepts c ON c.id = s.concept_id
    WHERE i.language_code = @language_code AND i.status = 'verified' AND NOT retired(s.id)
    ORDER BY e.id, s.position
) pool
ORDER BY random()
LIMIT 60;

-- name: CreateQualification :one
INSERT INTO qualification_attempts (user_id, language_code, expected)
VALUES (@user_id, @language_code, @expected)
RETURNING id;

-- name: LockQualification :one
SELECT id, user_id, language_code, expected, passed, answered_at, confirmed_at
FROM qualification_attempts WHERE id = $1 FOR UPDATE;

-- name: ScoreQualification :exec
UPDATE qualification_attempts SET score = @score::int, passed = @passed::boolean, answered_at = now()
WHERE id = @id;

-- name: ConfirmQualification :exec
UPDATE qualification_attempts SET confirmed_by = @confirmed_by, confirmed_at = now() WHERE id = @id;

-- name: GrantReviewer :exec
-- Approves a qualified user as a reviewer; an existing approval stays as it is.
INSERT INTO reviewer_approvals (user_id, language_code, level, approved_by)
VALUES (@user_id, @language_code, 'reviewer', @approved_by)
ON CONFLICT (user_id, language_code) DO NOTHING;

-- name: PendingQualifications :many
-- Passed attempts waiting for a senior reviewer of the languages, oldest first.
SELECT q.id, q.user_id, u.name, q.language_code, q.score, cardinality(q.expected)::int AS total, q.answered_at
FROM qualification_attempts q JOIN users u ON u.id = q.user_id
WHERE q.language_code = ANY(@languages::text[]) AND q.passed AND q.confirmed_at IS NULL
ORDER BY q.answered_at, q.id;

-- name: AddKeyUsage :exec
-- Adds one flush of counts, a row per key and day.
INSERT INTO key_usage (key_id, day, calls, errors, limited)
SELECT unnest(@key_ids::uuid[]), unnest(@days::date[]), unnest(@calls::int[]), unnest(@errors::int[]), unnest(@limited::int[])
ON CONFLICT (key_id, day) DO UPDATE SET
    calls = key_usage.calls + EXCLUDED.calls,
    errors = key_usage.errors + EXCLUDED.errors,
    limited = key_usage.limited + EXCLUDED.limited;

-- name: KeyUsage :many
-- A key's days with calls since a day, newest first. Check ownership with OwnKey.
SELECT day, calls, errors, limited FROM key_usage
WHERE key_id = @key_id AND day >= @since::date
ORDER BY day DESC;

-- name: OwnKey :one
SELECT name, prefix FROM api_keys WHERE id = @id AND user_id = @user_id;
