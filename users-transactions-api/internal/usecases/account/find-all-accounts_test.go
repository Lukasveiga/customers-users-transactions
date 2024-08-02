package usecases

import (
	"errors"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindAllAccountsUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	tenantId := int32(1)
	accounts := []domain.Account{
		{
			Id:       1,
			TenantId: 1,
			Number:   uuid.New().String(),
			Status:   "active",
		},
	}

	sut := NewFindAllAccountsUsecase(mockRepo)

	t.Run("Error to find all accounts", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindAll").Return(nil, expectedErr)
		defer mockRepo.On("FindAll").Unset()

		result, err := sut.FindAll(tenantId)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success find all accounts", func(t *testing.T) {
		mockRepo.On("FindAll").Return(accounts, nil)
		defer mockRepo.On("FindAll").Unset()

		result, err := sut.FindAll(tenantId)

		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, accounts[0], result[0])
	})
}
