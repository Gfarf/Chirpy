-- name: CreateStoredChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirpsAll :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetOneChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id = $1;

-- name: EditChirpByID :exec
UPDATE chirps SET body = $1, updated_at = NOW() WHERE id = $2;