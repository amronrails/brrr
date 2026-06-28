-- name: CreateCategory :one
INSERT INTO categories (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = $1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCategories :one
SELECT count(*) FROM categories;

-- name: UpdateCategory :one
UPDATE categories SET name = $2, slug = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
