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

func TestCreateAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Number:   uuid.New().String(),
		Status:   "active",
	}

	sut := NewCreateAccountUsecase(mockRepo)

	t.Run("Error to find accout", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindByNumber").Return(nil, expectedErr)
		defer mockRepo.On("FindByNumber").Unset()

		result, err := sut.Create(account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account number already exists", func(t *testing.T) {
		mockRepo.On("FindByNumber").Return(account, nil)
		defer mockRepo.On("FindByNumber").Unset()

		expectedErr := &shared.AlreadyExistsError{
			Object: "account",
			Id:     account.Number,
		}

		result, err := sut.Create(account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error to create account", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		mockRepo.On("Create").Return(nil, expectedErr)
		defer mockRepo.On("Create").Unset()

		result, err := sut.Create(account)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error invalid input to create account", func(t *testing.T) {
		invalidAccount := &domain.Account{
			Id:       1,
			TenantId: 1,
			Number:   "invalid",
			Status:   "invalid",
		}

		expectedError := &shared.ValidationError{
			Errors: map[string]string{
				"number": "must be a valid uuid",
				"status": "must be active or inactive",
			},
		}

		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		result, err := sut.Create(invalidAccount)

		assert.Nil(t, result)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("Success create new account", func(t *testing.T) {
		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		mockRepo.On("Create").Return(account, nil)
		defer mockRepo.On("Create").Unset()

		result, err := sut.Create(account)

		assert.Nil(t, err)
		assert.Equal(t, account, result)
	})
}
