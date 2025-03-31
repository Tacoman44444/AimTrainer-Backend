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

-- name: GetTopTenScores :many
SELECT * FROM sessions
INNER JOIN users
ON sessions.player_id = users.id
ORDER BY score DESC, sessions.created_at DESC
LIMIT 10;