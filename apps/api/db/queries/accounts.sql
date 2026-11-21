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
