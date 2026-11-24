-- name: NextReviewCard :many
-- The next card: the best-priority concept with senses, or phrases, still
-- open to this reviewer, and all of those items in random order (AI-3, AI-4,
-- CON-3). A card never mixes senses and phrases.
WITH open AS (
    SELECT a.item_id, a.concept_id, a.kind, a.priority, i.created_at
    FROM review_agreement a
    JOIN items i ON i.id = a.item_id
    WHERE a.language_code = @language_code::text
      AND open_for_review(i.status) AND NOT retired(i.id)
      AND i.author_id IS DISTINCT FROM @reviewer_id::text
      AND NOT EXISTS (SELECT 1 FROM reviews r WHERE r.item_id = i.id AND r.reviewer_id = @reviewer_id::text)
), next AS (
    SELECT o.concept_id, o.kind, min(o.priority) AS priority
    FROM open o
    WHERE o.concept_id <> ALL (coalesce(@skip::uuid[], '{}'))
    GROUP BY o.concept_id, o.kind
    ORDER BY min(o.priority), min(o.created_at), o.concept_id, o.kind
    LIMIT 1
)
SELECT o.item_id::uuid AS item_id, n.concept_id::uuid AS concept_id, n.kind::item_kind AS kind,
       n.priority::review_priority AS priority
FROM open o JOIN next n ON n.concept_id = o.concept_id AND n.kind = o.kind
ORDER BY random();

-- name: CreateReview :one
INSERT INTO reviews (item_id, reviewer_id, decision, comment, correction)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ListOwnReviews :many
-- The reviewer's reviews among ids; others' are left out.
SELECT id, item_id, created_at FROM reviews
WHERE id = ANY(@ids::uuid[]) AND reviewer_id = @reviewer_id;

-- name: DeleteReviews :exec
DELETE FROM reviews WHERE id = ANY(@ids::uuid[]);


-- name: CreateResolution :one
INSERT INTO resolutions (item_id, senior_id, outcome, note) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: SetReviewsAside :exec
UPDATE reviews SET set_aside_by = @resolution_id WHERE item_id = @item_id AND set_aside_by IS NULL;

-- name: GetLatestResolution :one
-- The resolution that counts for promotion, if any.
SELECT outcome FROM resolutions WHERE item_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: ListDisputes :many
-- Disputed senses and phrases in the languages, oldest first. A dispute is as old as what
-- made it: the last counted review, or the confirmation of a report.
SELECT i.id, i.language_code,
       greatest(
           (SELECT max(r.created_at) FROM reviews r WHERE r.item_id = i.id AND r.set_aside_by IS NULL),
           (SELECT max(p.handled_at) FROM reports p
            WHERE p.item_id = i.id AND p.state = 'confirmed' AND p.settled_by IS NULL)
       )::timestamptz AS disputed_since
FROM items i
WHERE i.status = 'disputed' AND i.kind IN ('sense', 'phrase', 'example')
  AND i.language_code = ANY(@languages::text[])
ORDER BY disputed_since, i.id
LIMIT 50;

-- name: ListCountedReviews :many
-- The reviews of the items that count: those no send-back set aside.
SELECT item_id, decision, comment, correction, created_at
FROM reviews
WHERE item_id = ANY(@ids::uuid[]) AND set_aside_by IS NULL
ORDER BY item_id, created_at;

-- name: HasReviewed :one
-- Whether the reviewer ever reviewed the item, set aside or not.
SELECT EXISTS (SELECT 1 FROM reviews WHERE item_id = $1 AND reviewer_id = $2);

-- name: ListProviderQuality :many
SELECT provider, language_code, decided, accepted, corrected
FROM provider_quality ORDER BY language_code, provider;

-- name: ListCredits :many
SELECT name, words, reviews FROM credits
WHERE language_code = $1
ORDER BY words + reviews DESC, name, user_id
LIMIT 500;

-- name: HasUnsettledReport :one
-- A confirmed report no resolution has settled yet makes the item disputed.
SELECT EXISTS (SELECT 1 FROM reports
               WHERE item_id = $1 AND state = 'confirmed' AND settled_by IS NULL);

-- name: SettleReports :exec
UPDATE reports SET settled_by = @resolution_id, updated_at = now()
WHERE item_id = @item_id AND state = 'confirmed' AND settled_by IS NULL;

-- name: CreateReport :one
-- Only a sense can be reported; no row means there is no such sense.
INSERT INTO reports (item_id, reason, reporter_id)
SELECT i.id, @reason, sqlc.narg(reporter_id) FROM items i WHERE i.id = @item_id AND i.kind = 'sense'
RETURNING id;

-- name: ListOpenReports :many
-- Open reports in the languages, oldest first, with the sense they are about.
SELECT p.id, p.item_id, p.reason, p.reporter_id, p.created_at, i.language_code, i.status,
       e.headword, s.gloss
FROM reports p
JOIN items i ON i.id = p.item_id
JOIN senses s ON s.id = p.item_id
JOIN entries e ON e.id = s.entry_id
WHERE p.state = 'open' AND i.language_code = ANY(@languages::text[])
ORDER BY p.created_at, p.id
LIMIT 100;

-- name: GetReportItem :one
SELECT item_id FROM reports WHERE id = $1;

-- name: LockReport :one
-- Taken after the item's lock, as resolutions take them in that order too.
SELECT reporter_id, state FROM reports WHERE id = $1 FOR UPDATE;

-- name: HandleReport :exec
UPDATE reports SET state = @state, handled_by = @handled_by, handled_at = now(), note = @note, updated_at = now()
WHERE id = @id;
