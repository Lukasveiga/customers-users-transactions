package mocks

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCardRepository struct {
	mock.Mock
}

func (mock *MockCardRepository) Save(card *domain.Card) (*domain.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Card), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockCardRepository) FindAllByAccountId(accountId int32) ([]*domain.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]*domain.Card), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockCardRepository) FindById(accountId int32, id int32) (*domain.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Card), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockCardRepository) Update(id int32, card *domain.Card) (*domain.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*domain.Card), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mock *MockCardRepository) Delete(id int32) error {
	args := mock.Called()
	return args.Error(0)
}
