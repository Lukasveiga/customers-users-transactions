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

func TestFindOneAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	tenantId := int32(1)
	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Number:   uuid.New().String(),
		Status:   "active",
	}

	sut := NewFindOneAccountUsecase(mockRepo)

	t.Run("Error to find account by id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(nil, expectedErr)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.FindOne(tenantId, account.Id)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error account not found", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, nil)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.FindOne(tenantId, account.Id)

		expectedErr := &shared.EntityNotFoundError{
			Object: "account",
			Id:     account.Id,
		}

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success find account by id", func(t *testing.T) {
		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.FindOne(tenantId, account.Id)

		assert.Nil(t, err)
		assert.Equal(t, result, account)
	})
}
