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

func TestFindOneTenantUseCase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)

	sut := NewFindOneTenantUseCase(mockRepo)

	t.Run("Error to find tenant by id", func(t *testing.T) {
		expectedErr := errors.New("internal repo error")

		mockRepo.On("GetTenant").Return(nil, expectedErr)
		defer mockRepo.On("GetTenant").Unset()

		result, err := sut.FindOne(1)

		assert.Empty(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Error tenant not found", func(t *testing.T) {
		mockRepo.On("GetTenant").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetTenant").Unset()

		result, err := sut.FindOne(1)

		expectedErr := &shared.EntityNotFoundError{
			Object: "tenant",
			Id:     1,
		}

		assert.Empty(t, result)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})

	t.Run("Success to find tenant by id", func(t *testing.T) {
		tenant := infra.Tenant{
			ID:   int32(1),
			Name: "Tenant A",
		}

		mockRepo.On("GetTenant").Return(tenant, nil)
		defer mockRepo.On("GetTenant").Unset()

		result, err := sut.FindOne(tenant.ID)

		assert.NoError(t, err)
		assert.Equal(t, *result, tenant)
	})
}
