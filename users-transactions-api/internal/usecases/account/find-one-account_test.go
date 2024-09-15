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

func TestFindOneAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	tenantId := int32(1)
	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	sut := NewFindOneAccountUsecase(mockRepo)

	t.Run("Error to find account by id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("GetAccount").Return(nil, expectedErr)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.FindOne(tenantId, account.ID)

		assert.Empty(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.FindOne(tenantId, account.ID)

		expectedErr := &shared.EntityNotFoundError{
			Object: "account",
			Id:     account.ID,
		}

		assert.Empty(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success find account by id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.FindOne(tenantId, account.ID)

		assert.Nil(t, err)
		assert.Equal(t, *result, account)
	})
}
