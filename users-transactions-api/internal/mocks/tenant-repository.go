package mocks

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTenantRepository struct {
	mock.Mock
}

func (mock *MockTenantRepository) FindById(id int32) (*domain.Tenant, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Tenant), args.Error(1)
	}

	return nil, args.Error(1)
}
