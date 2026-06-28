-- name: CreateComment :one
INSERT INTO comments (body, approved, post_id, author_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetComment :one
SELECT * FROM comments WHERE id = $1;

-- name: ListComments :many
SELECT * FROM comments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountComments :one
SELECT count(*) FROM comments;

-- name: UpdateComment :one
UPDATE comments SET body = $2, approved = $3, post_id = $4, author_id = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;
