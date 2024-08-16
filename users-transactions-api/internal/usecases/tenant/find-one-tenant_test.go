package usecases

import (
	"errors"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestFindOneTenantUseCase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockTenantRepository)

	sut := NewFindOneTenantUseCase(mockRepo)

	t.Run("Error to find tenant by id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("FindById").Return(nil, expectedErr)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.FindOne(1)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error tenant not found", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, nil)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.FindOne(1)

		expectedErr := &shared.EntityNotFoundError{
			Object: "tenant",
			Id:     1,
		}

		assert.Nil(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success to find tenant by id", func(t *testing.T) {
		tenant := &domain.Tenant{
			Id:   int32(1),
			Name: "Tenant A",
		}

		mockRepo.On("FindById").Return(tenant, nil)
		defer mockRepo.On("FindById").Unset()

		result, err := sut.FindOne(tenant.Id)

		assert.NoError(t, err)
		assert.Equal(t, result, tenant)
	})
}
