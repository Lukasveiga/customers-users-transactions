package mocks

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (mock *MockAccountRepository) Create(account *domain.Account) (*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Account), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockAccountRepository) FindAll(tenantId int32) ([]domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]domain.Account), args.Error(1)
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

func (mock *MockAccountRepository) FindByNumber(tenantId int32, number string) (*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Account), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockAccountRepository) Update(id int32, account *domain.Account) (*domain.Account, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.Account), args.Error(1)
}

func (mock *MockAccountRepository) Delete(id int32) error {
	args := mock.Called()
	return args.Error(1)
}
