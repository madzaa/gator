-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1 and updated_at = $1
WHERE id = $2
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at DESC NULLS FIRST;