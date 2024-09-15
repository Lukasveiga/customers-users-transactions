package usecases

import (
	"database/sql"
	"errors"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/mocks"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	sut := NewActiveAccountUsecase(mockRepo)

	t.Run("Error to find account by Id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("GetAccount").Return(nil, expectedErr)
		defer mockRepo.On("GetAccount").Unset()

		err := sut.Active(account.TenantID, account.ID)

		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account not found by id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		expectedErr := &shared.EntityNotFoundError{
			Object: "account",
			Id:     account.ID,
		}

		err := sut.Active(account.TenantID, account.ID)

		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error to active an account", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("UpdateAccount").Return(nil, expectedErr)
		defer mockRepo.On("UpdateAccount").Unset()

		err := sut.Active(account.TenantID, account.ID)

		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success to active an account", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("UpdateAccount").Return(account, nil)
		defer mockRepo.On("UpdateAccount").Unset()

		err := sut.Active(account.TenantID, account.ID)

		assert.NoError(t, err)
	})
}
