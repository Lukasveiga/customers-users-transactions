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

-- name: SearchTransactions :many
SELECT
t.id,
t.card_id,
t.kind,
t.value
FROM transactions t
JOIN cards c ON t.card_id = c.id
JOIN accounts a ON c.account_id = a.id
WHERE a.tenant_id = $1 AND a.id = sqlc.arg(accountId);
