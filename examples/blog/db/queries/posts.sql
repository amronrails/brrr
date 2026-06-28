-- name: CreatePost :one
INSERT INTO posts (title, slug, excerpt, body, published, views, author_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPosts :one
SELECT count(*) FROM posts;

-- name: UpdatePost :one
UPDATE posts SET title = $2, slug = $3, excerpt = $4, body = $5, published = $6, views = $7, author_id = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;
