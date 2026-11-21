-- Reviewer agreement rates (P1-REV-05a, DQ-7). A decision agrees when it
-- matches the item's settled outcome: accept for `verified`, reject or
-- correct for `rejected`. Items still open or disputed do not count yet, and
-- neither do reviews a send-back set aside, as for promotion.
--
-- Below the threshold a reviewer loses approval until retrained: the session
-- leaves such approvals out, so nothing has to run on a schedule. Marking a
-- reviewer retrained restarts their count.

-- +goose Up
ALTER TABLE reviewer_approvals ADD COLUMN retrained_at timestamptz;

CREATE VIEW reviewer_agreement AS
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

-- The threshold: under 80% agreement over at least 20 settled decisions.
-- +goose StatementBegin
CREATE FUNCTION below_agreement(decided int, agreed int) RETURNS boolean
LANGUAGE sql IMMUTABLE AS $$ SELECT decided >= 20 AND agreed < 0.8 * decided $$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION below_agreement;
DROP VIEW reviewer_agreement;
ALTER TABLE reviewer_approvals DROP COLUMN retrained_at;
