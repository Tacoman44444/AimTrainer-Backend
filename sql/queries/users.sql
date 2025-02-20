-- name: CreateUser :one
INSERT INTO users (id, username, hashed_password, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW()
) RETURNING *;
