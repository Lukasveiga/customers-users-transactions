package infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionTxRepository(t *testing.T) {
	transactionTx := NewTx(testDb)

	t.Run("[CreateTransactionTx] Concurrency", func(t *testing.T) {
		account := createTestAccount(t, 1)
		card := createTestCard(t, account.ID)

		assert.Equal(t, int64(0), card.Amount)

		n := 100
		value := int64(200)

		errs := make(chan error, n)
		results := make(chan Transaction, n)

		arg := CreateTransactionParams{
			CardID: card.ID,
			Kind:   "Streaming Z",
			Value:  value,
		}

		for i := 0; i < n; i++ {
			go func() {
				ctx := context.Background()
				result, err := transactionTx.CreateTransactionTx(ctx, arg)

				errs <- err
				results <- result
			}()
		}

		for i := 0; i < n; i++ {
			err := <-errs
			assert.NoError(t, err)

			result := <-results
			assert.NotEmpty(t, result)
			assert.Equal(t, arg.Kind, result.Kind)
			assert.Equal(t, arg.Value, result.Value)
			assert.NotEmpty(t, result.CreatedAt)
		}

		card2, err := transactionTx.GetCard(context.Background(), GetCardParams{
			AccountID: account.ID,
			ID:        card.ID,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(n*int(value)), card2.Amount)
	})
}
