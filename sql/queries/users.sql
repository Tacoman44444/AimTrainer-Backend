-- name: CreateUser :one
INSERT INTO users (id, username, hashed_password, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW()
) RETURNING *;

-- name: FindUserByUsername :one
SELECT id, hashed_password
FROM users
WHERE username = $1;