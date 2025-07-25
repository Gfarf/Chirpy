-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    false
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetHashedPasswordByEmail :one
SELECT hashed_password FROM users WHERE email = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users SET email = $1, hashed_password = $2, updated_at = NOW() WHERE id = $3;

-- name: UpgradeUser :exec
UPDATE users SET is_chirpy_red = true, updated_at = NOW() WHERE id = $1;
