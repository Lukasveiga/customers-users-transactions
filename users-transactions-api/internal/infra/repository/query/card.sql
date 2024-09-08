-- name: CreateCard :one
INSERT INTO cards (
    account_id
) VALUES (
    $1
) RETURNING *;

-- name: GetCard :one
SELECT * FROM cards 
WHERE account_id = $1 AND id = $2
LIMIT 1;

-- name: GetCards :many
SELECT * FROM cards 
WHERE account_id = $1;

-- name: AddAmount :one
UPDATE cards 
SET amount = amount + $2,
updated_at = $3
WHERE id = $1 RETURNING *;