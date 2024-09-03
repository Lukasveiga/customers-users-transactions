package infra

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestAccount(t *testing.T, tenandId int32) Account {
	arg := CreateAccountParams{
		TenantID: sql.NullInt32{
			Int32: tenandId,
			Valid: true,
		},
		Status: "active",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, account)
	assert.Equal(t, arg.TenantID, account.TenantID)
	assert.Equal(t, arg.Status, account.Status)

	return account
}

func TestAccountRepository(t *testing.T) {

	t.Run("[CreateAccount] should create new account and return it", func(t *testing.T) {
		createTestAccount(t, 1)
	})

	t.Run("[GetAccount] should find account by id and tenant id", func(t *testing.T) {
		account := createTestAccount(t, 1)
		foundAccount, err := testQueries.GetAccount(context.Background(), GetAccountParams{
			TenantID: account.TenantID,
			ID:       account.ID,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, foundAccount)
		assert.Equal(t, account, foundAccount)
	})

	t.Run("[GetAccounts] shoudl return a list of accounts given a tenant id", func(t *testing.T) {
		n := 5

		for i := 0; i < n; i++ {
			createTestAccount(t, 2)
		}

		accounts, err := testQueries.GetAccounts(context.Background(), sql.NullInt32{
			Int32: int32(2),
			Valid: true,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, accounts)
		assert.Len(t, accounts, n)

		for _, account := range accounts {
			assert.NotEmpty(t, account)
		}
	})

	t.Run("[UpdateAccount] should update account and return it", func(t *testing.T) {
		account := createTestAccount(t, 1)

		updateArgs := UpdateAccountParams{
			TenantID: account.TenantID,
			Status:   "inactive",
			ID:       account.ID,
		}

		updatedAccount, err := testQueries.UpdateAccount(context.Background(), updateArgs)

		assert.NoError(t, err)
		assert.NotEmpty(t, updatedAccount)
		assert.Equal(t, updateArgs.Status, updatedAccount.Status)
	})
}
