-- name: CreateProduct :one
INSERT INTO products (name, sku, price, stock, active, metadata, category_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT count(*) FROM products;

-- name: UpdateProduct :one
UPDATE products SET name = $2, sku = $3, price = $4, stock = $5, active = $6, metadata = $7, category_id = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
