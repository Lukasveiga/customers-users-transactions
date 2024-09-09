-- name: CreateTransaction :one
INSERT INTO transactions (
    card_id,
    kind,
    value
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions 
WHERE card_id = $1 AND id = $2
LIMIT 1;

-- name: GetTransactions :many
SELECT * FROM transactions 
WHERE card_id = $1;