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
INSERT INTO reviews (item_id, reviewer_id, decision, comment)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ListOwnReviews :many
-- The reviewer's reviews among ids; others' are left out.
SELECT id, item_id, created_at FROM reviews
WHERE id = ANY(@ids::uuid[]) AND reviewer_id = @reviewer_id;

-- name: DeleteReviews :exec
DELETE FROM reviews WHERE id = ANY(@ids::uuid[]);

-- name: ListItemDecisions :many
SELECT decision FROM reviews WHERE item_id = $1;
