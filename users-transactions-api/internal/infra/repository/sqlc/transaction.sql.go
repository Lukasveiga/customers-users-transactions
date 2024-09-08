// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transaction.sql

package infra

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    card_id,
    kind,
    value
) VALUES (
    $1, $2, $3
) RETURNING id, card_id, kind, value, created_at, updated_at, deleted_at
`

type CreateTransactionParams struct {
	CardID int32  `json:"card_id"`
	Kind   string `json:"kind"`
	Value  int64  `json:"value"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction, arg.CardID, arg.Kind, arg.Value)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.CardID,
		&i.Kind,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
