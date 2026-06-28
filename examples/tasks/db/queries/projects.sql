-- name: CreateProject :one
INSERT INTO projects (name, key, description, archived)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProjects :one
SELECT count(*) FROM projects;

-- name: UpdateProject :one
UPDATE projects SET name = $2, key = $3, description = $4, archived = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;
