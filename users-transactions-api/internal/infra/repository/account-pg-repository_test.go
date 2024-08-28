package infra

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgreAccountRepository(t *testing.T) {
	dbClient := config.SetupPgTestcontainers()

	pgAccountRepository := NewPgAccountRepository(dbClient)

	defer dbClient.Close()

	t.Run("should successfully connect to pg container", func(t *testing.T) {
		assert.NotNil(t, &dbClient)
	})

	t.Run("[Save] should save new account and return it", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 1,
			Status:   "active",
		}

		savedAccount, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		account.Id = savedAccount.Id
		account.CreatedAt = savedAccount.CreatedAt
		account.UpdatedAt = savedAccount.UpdatedAt
		account.DeletedAt = savedAccount.DeletedAt

		assert.Equal(t, account, savedAccount)
	})

	t.Run("[Save] should return error when tenant id not exists", func(t *testing.T) {
		invalidTenantId := int32(15)
		account := &domain.Account{
			TenantId: invalidTenantId,
			Status:   "active",
		}

		savedAccount, err := pgAccountRepository.Save(account)

		assert.Nil(t, savedAccount)
		assert.NotNil(t, err)
		assert.Contains(t, err.(*pq.Error).Detail,
			fmt.Sprintf("Key (tenant_id)=(%d) is not present in table \"tenants\".", invalidTenantId))
	})

	t.Run("[FindById] should find account by id and tenant id", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 1,
			Status:   "active",
		}

		savedAccount, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		foundAccount, err := pgAccountRepository.FindById(account.TenantId, savedAccount.Id)

		assert.NoError(t, err)

		assert.NotNil(t, foundAccount)
		assert.Equal(t, savedAccount, foundAccount)
	})

	t.Run("[FindById] should return nil when account is not found by id", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 1,
			Status:   "active",
		}

		savedAccount, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		testCases := []struct {
			name      string
			tenantId  int32
			accountId int32
		}{
			{
				name:      "invalid tenant id",
				tenantId:  15,
				accountId: savedAccount.Id,
			},
			{
				name:      "invalid account id",
				tenantId:  1,
				accountId: 5,
			},
		}

		for i := range testCases {
			tc := testCases[i]

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				foundAccount, err := pgAccountRepository.FindById(tc.tenantId, tc.accountId)

				assert.Nil(t, foundAccount)
				assert.Nil(t, err)
			})
		}
	})

	t.Run("[FindAll] should return a list of account given a tenant id", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 2,
			Status:   "active",
		}

		savedAccount, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		accounts, err := pgAccountRepository.FindAll(account.TenantId)

		assert.NoError(t, err)

		assert.Len(t, accounts, 1)
		assert.Equal(t, savedAccount, accounts[0])
	})

	t.Run("[FindAll] should return an empty list if it has no account related to a given tenant id",
		func(t *testing.T) {
			accounts, err := pgAccountRepository.FindAll(3)

			assert.NoError(t, err)
			assert.Len(t, accounts, 0)
		})

	t.Run("[Update] should update account", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 3,
			Status:   "active",
		}

		savedAccount, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		updateAccount := &domain.Account{
			Id:       savedAccount.Id,
			TenantId: 3,
			Status:   "inactive",
		}

		updatedAccount, err := pgAccountRepository.Update(updateAccount)

		assert.NoError(t, err)
		assert.Equal(t, updateAccount.Status, updatedAccount.Status)
	})

	t.Run("[Update] should return error when update account with invalid account id", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 3,
			Status:   "active",
		}

		account.Id = 15

		updatedAccount, err := pgAccountRepository.Update(account)

		assert.Nil(t, updatedAccount)
		assert.Equal(t, sql.ErrNoRows.Error(), err.Error())
	})

	t.Run("[Update] should return error when update account with invalid tenant id", func(t *testing.T) {
		account := &domain.Account{
			TenantId: 3,
			Status:   "active",
		}

		_, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		updateAccount := &domain.Account{
			TenantId: 15,
			Status:   "inactive",
		}

		updatedAccount, err := pgAccountRepository.Update(updateAccount)

		assert.Nil(t, updatedAccount)
		assert.Equal(t, sql.ErrNoRows.Error(), err.Error())
	})
}
