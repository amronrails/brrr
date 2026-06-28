-- name: CreateTask :one
INSERT INTO tasks (title, description, status, priority, due_date, done, project_id, assignee_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountTasks :one
SELECT count(*) FROM tasks;

-- name: UpdateTask :one
UPDATE tasks SET title = $2, description = $3, status = $4, priority = $5, due_date = $6, done = $7, project_id = $8, assignee_id = $9, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;
