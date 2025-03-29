-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2, false) RETURNING *;

-- name: DeleteUsers :exec
DELETE
FROM users;

-- name: UpdateUserEmailAndPassword :one
UPDATE users
set hashed_password = $1,
    email           = $2,
    updated_at      = NOW()
WHERE id = $3 RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
where email = $1;

-- name: UpdateUserToChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;