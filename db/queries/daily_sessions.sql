-- name: FindDailySessionById :one
SELECT * FROM daily_sessions
WHERE id = $1 LIMIT 1;

-- name: FindAllDailySessions :many
SELECT * FROM daily_sessions
ORDER BY id;

-- name: CreateDailySession :one
INSERT INTO daily_sessions (session_date, raw_notes_blob, generated_script)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteDailySession :exec
DELETE FROM daily_sessions
WHERE id = $1;
