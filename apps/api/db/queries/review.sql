-- name: NextReviewCard :many
-- The next card: the best-priority concept with senses still open to this
-- reviewer, and all of those senses in random order (AI-3, AI-4, CON-3).
WITH open AS (
    SELECT a.item_id, a.concept_id, a.priority, i.created_at
    FROM review_agreement a
    JOIN items i ON i.id = a.item_id
    WHERE a.language_code = @language_code::text
      AND open_for_review(i.status)
      AND i.author_id IS DISTINCT FROM @reviewer_id::text
      AND NOT EXISTS (SELECT 1 FROM reviews r WHERE r.item_id = i.id AND r.reviewer_id = @reviewer_id::text)
), next AS (
    SELECT o.concept_id, min(o.priority) AS priority
    FROM open o
    WHERE o.concept_id <> ALL (coalesce(@skip::uuid[], '{}'))
    GROUP BY o.concept_id
    ORDER BY min(o.priority), min(o.created_at), o.concept_id
    LIMIT 1
)
SELECT o.item_id::uuid AS item_id, n.concept_id::uuid AS concept_id, n.priority::review_priority AS priority
FROM open o JOIN next n ON n.concept_id = o.concept_id
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
-- Disputed senses in the languages, oldest first. A dispute is as old as the
-- last counted review, which is the one that made the reviewers disagree.
SELECT i.id, i.language_code, max(r.created_at)::timestamptz AS disputed_since
FROM items i
JOIN reviews r ON r.item_id = i.id AND r.set_aside_by IS NULL
WHERE i.status = 'disputed' AND i.kind = 'sense'
  AND i.language_code = ANY(@languages::text[])
GROUP BY i.id
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
