-- name: CreateAccount :one
INSERT INTO accounts (
    tenant_id, status
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts 
WHERE tenant_id = $1 AND id = $2
LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts 
WHERE tenant_id = $1;

-- name: UpdateAccount :one
UPDATE accounts 
SET status = $2, 
updated_at = $3, 
deleted_at = $4 
WHERE id = $1
RETURNING *;