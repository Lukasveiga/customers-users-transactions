package usecases

import (
	"errors"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Number:   uuid.New().String(),
		Status:   "active",
	}

	sut := NewDeleteAccountUsecase(mockRepo)

	t.Run("Error to find account by Id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(nil, expectedErr)
		defer mockRepo.On("FindById").Unset()

		err := sut.Delete(account.TenantId, account.Id)

		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account not found by id", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, nil)
		defer mockRepo.On("FindById").Unset()

		expectedErr := &shared.EntityNotFoundError{
			Object: "account",
			Id:     account.Id,
		}

		err := sut.Delete(account.TenantId, account.Id)

		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error to delete account", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("Delete").Return(expectedErr)
		defer mockRepo.On("Delete").Unset()

		err := sut.Delete(account.TenantId, account.Id)

		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success to delete account", func(t *testing.T) {
		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("Delete").Return(nil)
		defer mockRepo.On("Delete").Unset()

		err := sut.Delete(account.TenantId, account.Id)

		assert.NoError(t, err)
	})
}
