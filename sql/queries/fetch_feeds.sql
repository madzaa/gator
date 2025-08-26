-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1,  updated_at = $1
WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at DESC NULLS FIRST;