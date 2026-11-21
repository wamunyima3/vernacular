-- Public credit for volunteers, per language (CON-5, P1-REV-07): contributors
-- for their words that were verified, reviewers for decisions that matched a
-- settled outcome. Only people who chose to be credited are listed, and when
-- they chose is recorded (NFR-6). Banned users are never listed.

-- +goose Up
ALTER TABLE users ADD COLUMN credit_consent_at timestamptz;

-- agrees is the one definition of a decision matching a settled outcome,
-- shared by agreement rates and credit.
-- +goose StatementBegin
CREATE FUNCTION agrees(decision review_decision, status trust_status) RETURNS boolean
LANGUAGE sql IMMUTABLE AS $$ SELECT (decision = 'accept') = (status = 'verified') $$;
-- +goose StatementEnd

CREATE OR REPLACE VIEW reviewer_agreement AS
SELECT a.user_id, a.language_code,
       count(i.id)::int AS decided,
       count(i.id) FILTER (WHERE agrees(r.decision, i.status))::int AS agreed
FROM reviewer_approvals a
LEFT JOIN reviews r ON r.reviewer_id = a.user_id
    AND r.created_at > coalesce(a.retrained_at, '-infinity')
    AND r.set_aside_by IS NULL
LEFT JOIN items i ON i.id = r.item_id
    AND i.language_code = a.language_code
    AND i.status IN ('verified', 'rejected')
GROUP BY a.user_id, a.language_code;

CREATE VIEW credits AS
WITH words AS (
    SELECT author_id AS user_id, language_code, count(*)::int AS n
    FROM items
    WHERE kind = 'sense' AND status = 'verified' AND author_id IS NOT NULL
    GROUP BY author_id, language_code
), reviewed AS (
    SELECT r.reviewer_id AS user_id, i.language_code, count(*)::int AS n
    FROM reviews r
    JOIN items i ON i.id = r.item_id AND i.status IN ('verified', 'rejected')
    WHERE r.set_aside_by IS NULL AND agrees(r.decision, i.status)
    GROUP BY r.reviewer_id, i.language_code
)
SELECT u.id AS user_id, u.name, c.language_code,
       coalesce(w.n, 0) AS words, coalesce(v.n, 0) AS reviews
FROM (SELECT user_id, language_code FROM words UNION SELECT user_id, language_code FROM reviewed) c
JOIN users u ON u.id = c.user_id
LEFT JOIN words w ON w.user_id = c.user_id AND w.language_code = c.language_code
LEFT JOIN reviewed v ON v.user_id = c.user_id AND v.language_code = c.language_code
WHERE u.credit_consent_at IS NOT NULL AND NOT coalesce(u.banned, false);

-- +goose Down
DROP VIEW credits;
CREATE OR REPLACE VIEW reviewer_agreement AS
SELECT a.user_id, a.language_code,
       count(i.id)::int AS decided,
       count(i.id) FILTER (WHERE (r.decision = 'accept') = (i.status = 'verified'))::int AS agreed
FROM reviewer_approvals a
LEFT JOIN reviews r ON r.reviewer_id = a.user_id
    AND r.created_at > coalesce(a.retrained_at, '-infinity')
    AND r.set_aside_by IS NULL
LEFT JOIN items i ON i.id = r.item_id
    AND i.language_code = a.language_code
    AND i.status IN ('verified', 'rejected')
GROUP BY a.user_id, a.language_code;
DROP FUNCTION agrees;
ALTER TABLE users DROP COLUMN credit_consent_at;
