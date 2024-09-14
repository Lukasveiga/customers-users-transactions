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

	t.Run("[GetTransaction] should get transaction by id", func(t *testing.T) {
		transaction := createTestTransaction(t, 1)

		transaction2, err := testQueries.GetTransaction(context.Background(), GetTransactionParams{
			CardID: transaction.CardID,
			ID:     transaction.ID,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, transaction2)
		assert.Equal(t, transaction, transaction2)
	})

	t.Run("[GetTransactions] should get transactions by card id", func(t *testing.T) {
		transaction := createTestTransaction(t, 1)

		transactions, err := testQueries.GetTransactions(context.Background(), transaction.CardID)

		assert.NoError(t, err)
		assert.NotEmpty(t, transactions)
		assert.Len(t, transactions, 1)

		for _, transaction2 := range transactions {
			assert.NotEmpty(t, transaction2)
		}
	})

	t.Run("[SearchTransactions] should return filtered transactions", func(t *testing.T) {
		tenantId := int32(2)
		account := createTestAccount(t, tenantId)
		card := createTestCard(t, account.ID)
		ctx := context.Background()

		inputParams := CreateTransactionParams{
			CardID: card.ID,
			Kind:   "Streaming Z",
			Value:  100,
		}

		for i := 0; i < 5; i++ {
			testQueries.CreateTransaction(ctx, inputParams)
		}

		result, err := testQueries.SearchTransactions(ctx, SearchTransactionsParams{
			TenantID:  tenantId,
			Accountid: account.ID,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 5)

		for _, trans := range result {
			assert.NotEmpty(t, trans)
			assert.Equal(t, inputParams.Kind, trans.Kind)
			assert.Equal(t, inputParams.Value, trans.Value)
		}
	})
}
