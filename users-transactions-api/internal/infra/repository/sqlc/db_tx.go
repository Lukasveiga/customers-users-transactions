package infra

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type QuerierTx interface {
	CreateTransactionTx(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
}

type Tx struct {
	*Queries
	db *sql.DB
}

func NewTx(db *sql.DB) *Tx {
	return &Tx{
		db:      db,
		Queries: New(db),
	}
}

func (tx *Tx) CreateTransactionTx(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	var transaction Transaction
	var err error

	err = tx.execTx(ctx, func(q *Queries) error {

		transaction, err = q.CreateTransaction(ctx, arg)

		if err != nil {
			return err
		}

		_, err = q.AddAmount(ctx, AddAmountParams{
			ID:     arg.CardID,
			Amount: arg.Value,
			UpdatedAt: sql.NullTime{
				Time: time.Now().UTC(),
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

	return transaction, err
}

func (transactionTx *Tx) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := transactionTx.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			slog.Error(
				"error transaction tx",
				slog.String("err", err.Error()),
				slog.String("errTx", rbErr.Error()),
			)
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}
