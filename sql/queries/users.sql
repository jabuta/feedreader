-- name: CreateUser :one
INSERT INTO users (id, name, api_key, created_at, updated_at)
VALUES ($1, $2, encode(sha256(random()::text::bytea), 'hex'), NOW(), NOW())
RETURNING *;

-- name: GetUserByAPI :one
SELECT * FROM users
WHERE api_key = $1;