package usecases

import (
	"errors"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindAllAccountsUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	tenantId := int32(1)
	accounts := []infra.Account{
		{
			ID:       1,
			TenantID: 1,
			Status:   "active",
		},
	}

	sut := NewFindAllAccountsUsecase(mockRepo)

	t.Run("Error to find all accounts", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("GetAccounts").Return(nil, expectedErr)
		defer mockRepo.On("GetAccounts").Unset()

		result, err := sut.FindAll(tenantId)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success find all accounts", func(t *testing.T) {
		mockRepo.On("GetAccounts").Return(accounts, nil)
		defer mockRepo.On("GetAccounts").Unset()

		result, err := sut.FindAll(tenantId)

		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, accounts[0], result[0])
	})
}
