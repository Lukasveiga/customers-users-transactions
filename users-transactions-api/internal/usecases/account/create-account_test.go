package usecases

import (
	"errors"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Status:   "active",
	}

	sut := NewCreateAccountUsecase(mockRepo)

	t.Run("Error to create account", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("Save").Return(nil, expectedErr)
		defer mockRepo.On("Save").Unset()

		result, err := sut.Create(1)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success create new account", func(t *testing.T) {
		mockRepo.On("Save").Return(account, nil)
		defer mockRepo.On("Save").Unset()

		result, err := sut.Create(1)

		assert.Nil(t, err)
		assert.Equal(t, account, result)
	})
}
