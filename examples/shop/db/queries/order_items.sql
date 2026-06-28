-- name: CreateOrderItem :one
INSERT INTO order_items (quantity, unit_price, order_id, product_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_items WHERE id = $1;

-- name: ListOrderItems :many
SELECT * FROM order_items
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountOrderItems :one
SELECT count(*) FROM order_items;

-- name: UpdateOrderItem :one
UPDATE order_items SET quantity = $2, unit_price = $3, order_id = $4, product_id = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items WHERE id = $1;
