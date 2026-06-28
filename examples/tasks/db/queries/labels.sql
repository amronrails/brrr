-- name: CreateLabel :one
INSERT INTO labels (name, color)
VALUES ($1, $2)
RETURNING *;

-- name: GetLabel :one
SELECT * FROM labels WHERE id = $1;

-- name: ListLabels :many
SELECT * FROM labels
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountLabels :one
SELECT count(*) FROM labels;

-- name: UpdateLabel :one
UPDATE labels SET name = $2, color = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteLabel :exec
DELETE FROM labels WHERE id = $1;
