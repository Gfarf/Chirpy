-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: ValidateToken :one
SELECT * FROM refresh_tokens WHERE id = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens SET updated_at = NOW(), revoked_at = NOW() WHERE id = $1;