-- name: CreateAccount :one
INSERT INTO accounts (
    tenant_id, 
    status, 
    created_at, 
    updated_at, 
    deleted_at
) VALUES (
    $1, $2, $3, $4, $5
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
SET tenant_id = $1, 
status = $2, 
created_at = $3, 
updated_at = $4, 
deleted_at = $5 
WHERE id = $6
RETURNING *;