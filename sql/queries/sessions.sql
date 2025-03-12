-- name: CreateSession :one
INSERT INTO sessions (id, score, accuracy, created_at, player_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    $3
) RETURNING *;

-- name: GetPlayerBestSession :one
SELECT * FROM sessions
WHERE player_id = $1
ORDER BY score DESC
LIMIT 1;
