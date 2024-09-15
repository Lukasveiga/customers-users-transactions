package mocks

import (
	"context"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// Tenant
func (mock *MockRepository) GetTenant(ctx context.Context, id int32) (infra.Tenant, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Tenant), args.Error(1)
	}

	return infra.Tenant{}, args.Error(1)
}

// Account
func (mock *MockRepository) CreateAccount(ctx context.Context, arg infra.CreateAccountParams) (infra.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Account), args.Error(1)
	}

	return infra.Account{}, args.Error(1)
}

func (mock *MockRepository) GetAccount(ctx context.Context, arg infra.GetAccountParams) (infra.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Account), args.Error(1)
	}

	return infra.Account{}, args.Error(1)
}

func (mock *MockRepository) GetAccounts(ctx context.Context, tenantID int32) ([]infra.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]infra.Account), args.Error(1)
	}

	return []infra.Account{}, args.Error(1)
}

func (mock *MockRepository) UpdateAccount(ctx context.Context, arg infra.UpdateAccountParams) (infra.Account, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Account), args.Error(1)
	}

	return infra.Account{}, args.Error(1)
}

// Card
func (mock *MockRepository) CreateCard(ctx context.Context, accountId int32) (infra.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Card), args.Error(1)
	}

	return infra.Card{}, args.Error(1)
}

func (mock *MockRepository) GetCard(ctx context.Context, arg infra.GetCardParams) (infra.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Card), args.Error(1)
	}

	return infra.Card{}, args.Error(1)
}

func (mock *MockRepository) GetCards(ctx context.Context, accountID int32) ([]infra.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]infra.Card), args.Error(1)
	}

	return []infra.Card{}, args.Error(1)
}

func (mock *MockRepository) AddAmount(ctx context.Context, arg infra.AddAmountParams) (infra.Card, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Card), args.Error(1)
	}

	return infra.Card{}, args.Error(1)
}

// Transaction
func (mock *MockRepository) CreateTransaction(ctx context.Context, arg infra.CreateTransactionParams) (infra.Transaction, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Transaction), args.Error(1)
	}

	return infra.Transaction{}, args.Error(1)
}

func (mock *MockRepository) CreateTransactionTx(ctx context.Context, arg infra.CreateTransactionParams) (infra.Transaction, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Transaction), args.Error(1)
	}

	return infra.Transaction{}, args.Error(1)
}

func (mock *MockRepository) GetTransaction(ctx context.Context, arg infra.GetTransactionParams) (infra.Transaction, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.(infra.Transaction), args.Error(1)
	}

	return infra.Transaction{}, args.Error(1)
}

func (mock *MockRepository) GetTransactions(ctx context.Context, cardID int32) ([]infra.Transaction, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]infra.Transaction), args.Error(1)
	}

	return []infra.Transaction{}, args.Error(1)
}

func (mock *MockRepository) SearchTransactions(ctx context.Context, arg infra.SearchTransactionsParams) ([]infra.SearchTransactionsRow, error) {
	args := mock.Called()
	result := args.Get(0)

	if result != nil {
		return result.([]infra.SearchTransactionsRow), args.Error(1)
	}

	return []infra.SearchTransactionsRow{}, args.Error(1)
}
