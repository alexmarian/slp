-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1,
        NOW(),
        NOW(),
        $2,
        NOW() + INTERVAL '60 days') RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
    set revoked_at = NOW(),
    updated_at = NOW()
where token = $1;

-- name: GetValidRefreshToken :one
SELECT *  FROM refresh_tokens where token=$1 AND revoked_at IS NULL AND expires_at > NOW();