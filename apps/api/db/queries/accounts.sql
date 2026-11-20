-- name: GetSessionUser :one
-- The user behind a live session token. Banned users have no valid session.
SELECT u.id, u.name, u.email, coalesce(u.role, 'developer')::text AS role
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token = $1
  AND s.expires_at > now()
  AND NOT (coalesce(u.banned, false) AND (u.ban_expires IS NULL OR u.ban_expires > now()));

-- name: ListApprovals :many
SELECT language_code, level
FROM reviewer_approvals
WHERE user_id = $1
ORDER BY language_code;
