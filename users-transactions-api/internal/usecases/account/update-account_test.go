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

func TestUpdateAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Number:   uuid.New().String(),
		Status:   "active",
	}

	sut := NewUpdateAccountUsecase(mockRepo)

	t.Run("Error to find account by Id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(nil, expectedErr)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.Update(account.Id, account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account not found by id", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, nil)
		defer mockRepo.On("FindById").Unset()

		expectedErr := &shared.EntityNotFoundError{
			Object: "account",
			Id:     account.Id,
		}

		result, err := sut.Update(account.Id, account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error to find account by number", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("FindByNumber").Return(nil, expectedErr)
		defer mockRepo.On("FindByNumber").Unset()

		result, err := sut.Update(account.Id, account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account number already exists", func(t *testing.T) {
		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("FindByNumber").Return(account, nil)
		defer mockRepo.On("FindByNumber").Unset()

		expectedErr := &shared.AlreadyExistsError{
			Object: "account",
			Id:     account.Number,
		}

		updateAccount := &domain.Account{
			Id:       2,
			TenantId: 1,
			Number:   account.Number,
			Status:   "active",
		}

		result, err := sut.Update(updateAccount.Id, updateAccount)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error invalid input", func(t *testing.T) {
		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		expectedErr := &shared.ValidationError{
			Errors: map[string]string{
				"number": "must be a valid uuid",
				"status": "must be active or inactive",
			},
		}

		updateAccount := &domain.Account{
			Id:       2,
			TenantId: 1,
			Number:   "invalid",
			Status:   "invalid",
		}

		result, err := sut.Update(updateAccount.Id, updateAccount)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error to update account", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		mockRepo.On("Update").Return(nil, expectedErr)
		defer mockRepo.On("Update").Unset()

		result, err := sut.Update(account.Id, account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success to update account", func(t *testing.T) {
		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		mockRepo.On("Update").Return(account, nil)
		defer mockRepo.On("Update").Unset()

		result, err := sut.Update(account.Id, account)

		assert.Nil(t, err)
		assert.Equal(t, account, result)
	})
}
