-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.id AS user_id, refresh_tokens.*
FROM refresh_tokens
JOIN users ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $3
WHERE token = $1;