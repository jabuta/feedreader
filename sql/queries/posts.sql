-- name: CreatePost :one
INSERT INTO posts(id, title, url, description, published_at, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
INNER JOIN feed_follows
ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
LIMIT $2;