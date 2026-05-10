-- name: FindSubTaskById :one
SELECT * FROM sub_tasks
WHERE id = $1 LIMIT 1;

-- name: FindAllSubTasks :many
SELECT * FROM sub_tasks
ORDER BY id;

-- name: FindSubTasksByTaskId :many
SELECT * FROM sub_tasks
WHERE task_id = $1
ORDER BY id;

-- name: CreateSubTask :one
INSERT INTO sub_tasks (task_id, title, is_completed)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSubTask :one
UPDATE sub_tasks
SET title = $2, is_completed = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteSubTask :exec
DELETE FROM sub_tasks
WHERE id = $1;
