-- name: FindNoteById :one
SELECT * FROM task_notes
WHERE id = $1 LIMIT 1;

-- name: FindAllNotes :many
SELECT * FROM task_notes
ORDER BY id;

-- name: FindNotesByTaskId :many
SELECT * FROM task_notes
WHERE task_id = $1
ORDER BY id;

-- name: CreateNote :one
INSERT INTO task_notes (task_id, sub_task_id, content, category)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateNote :one
UPDATE task_notes
SET task_id = $2, sub_task_id = $3, content = $4, category = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM task_notes
WHERE id = $1;
