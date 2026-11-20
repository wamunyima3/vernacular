-- name: GetSessionUser :one
-- The user behind a live session token, with their review approvals, in one
-- round trip. Banned users have no valid session.
SELECT u.id, u.name, u.email, u.role,
       coalesce(
           json_agg(json_build_object('language', a.language_code, 'level', a.level) ORDER BY a.language_code)
               FILTER (WHERE a.user_id IS NOT NULL),
           '[]'
       )::jsonb AS approvals
FROM sessions s
JOIN users u ON u.id = s.user_id
LEFT JOIN reviewer_approvals a ON a.user_id = u.id
WHERE s.token = $1
  AND s.expires_at > now()
  AND NOT (coalesce(u.banned, false) AND (u.ban_expires IS NULL OR u.ban_expires > now()))
GROUP BY u.id;
