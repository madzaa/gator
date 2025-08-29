-- name: CreatePost :exec
INSERT INTO posts(id,
                  created_at,
                  updated_at,
                  title,
                  url,
                  description,
                  published_at,
                  feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.title as posts
FROM posts
         LEFT JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE user_id = $1
ORDER BY published_at DESC
LIMIT $2;

-- name: GetPost :one
SELECT *
FROM posts
WHERE url = $1;