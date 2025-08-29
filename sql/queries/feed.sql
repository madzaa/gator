-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: DeleteFeeds :exec
DELETE
FROM feeds;

-- name: GetFeeds :many
SELECT u.name as username, f.name as feedname
from feeds f
         left join users u on f.user_id = u.id;