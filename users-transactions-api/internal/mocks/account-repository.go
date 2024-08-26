package mocks

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (mock *MockAccountRepository) Save(account *domain.Account) (*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Account), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockAccountRepository) FindAll(tenantId int32) ([]*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]*domain.Account), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockAccountRepository) FindById(tenantId int32, id int32) (*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Account), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockAccountRepository) Update(account *domain.Account) (*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Account), args.Error(1)
	}

	return nil, args.Error(1)
}
