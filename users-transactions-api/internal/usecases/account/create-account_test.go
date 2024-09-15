package usecases

import (
	"errors"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	sut := NewCreateAccountUsecase(mockRepo)

	t.Run("Error to create account", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("CreateAccount").Return(nil, expectedErr)
		defer mockRepo.On("CreateAccount").Unset()

		result, err := sut.Create(1)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success create new account", func(t *testing.T) {
		mockRepo.On("CreateAccount").Return(account, nil)
		defer mockRepo.On("CreateAccount").Unset()

		result, err := sut.Create(1)

		assert.Nil(t, err)
		assert.Equal(t, account, *result)
	})
}
