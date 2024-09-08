-- name: CreateTransaction :one
INSERT INTO transactions (
    card_id,
    kind,
    value
) VALUES (
    $1, $2, $3
) RETURNING *;