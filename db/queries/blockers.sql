-- name: FindBlockerById :one
SELECT * FROM blockers
WHERE id = $1 LIMIT 1;

-- name: FindAllBlockers :many
SELECT * FROM blockers
ORDER BY id;

-- name: FindBlockersByTaskId :many
SELECT * FROM blockers
WHERE task_id = $1
ORDER BY id;

-- name: FindBlockersByTaskIdAndSubTaskId :many
SELECT * FROM blockers
WHERE task_id = $1 AND (sub_task_id = $2 OR ($2 IS NULL AND sub_task_id IS NULL))
ORDER BY id;

-- name: CreateBlocker :one
INSERT INTO blockers (task_id, sub_task_id, reason, is_resolved)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateBlocker :one
UPDATE blockers
SET task_id = $2, sub_task_id = $3, reason = $4, is_resolved = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteBlocker :exec
DELETE FROM blockers
WHERE id = $1;
