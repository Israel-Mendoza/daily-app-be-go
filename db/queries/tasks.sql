-- name: FindTaskById :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: FindAllTasks :many
SELECT * FROM tasks
ORDER BY id;

-- name: CreateTask :one
INSERT INTO tasks (title, status)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2, status = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
