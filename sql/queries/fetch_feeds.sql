-- name: MarkFeedFetched :exec
UPDATE feed_follows
SET last_fetched_at = $1 and updated_at = $1
WHERE id = $2
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT f.*
FROM feed_follows left join public.feeds f on feed_follows.feed_id = f.id
ORDER BY last_fetched_at DESC NULLS FIRST;