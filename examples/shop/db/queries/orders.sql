-- name: CreateOrder :one
INSERT INTO orders (status, total, placed_at, customer_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountOrders :one
SELECT count(*) FROM orders;

-- name: UpdateOrder :one
UPDATE orders SET status = $2, total = $3, placed_at = $4, customer_id = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;
