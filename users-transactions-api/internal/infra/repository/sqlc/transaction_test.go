package infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTransaction(t *testing.T, tenantId int32) Transaction {
	account := createTestAccount(t, tenantId)
	card := createTestCard(t, account.ID)

	arg := CreateTransactionParams{
		CardID: card.ID,
		Kind:   "Streamin Z",
		Value:  42,
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, transaction)
	assert.Equal(t, arg.Kind, transaction.Kind)
	assert.Equal(t, arg.Value, transaction.Value)
	assert.NotEmpty(t, transaction.CreatedAt)

	return transaction
}

func TestTransactionRepository(t *testing.T) {

	t.Run("[CreateTransaction] should create new transaction and return it", func(t *testing.T) {
		createTestTransaction(t, 1)
	})
}
