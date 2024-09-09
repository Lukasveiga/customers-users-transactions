// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package infra

import (
	"context"
)

type Querier interface {
	AddAmount(ctx context.Context, arg AddAmountParams) (Card, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateCard(ctx context.Context, accountID int32) (Card, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	GetAccount(ctx context.Context, arg GetAccountParams) (Account, error)
	GetAccounts(ctx context.Context, tenantID int32) ([]Account, error)
	GetCard(ctx context.Context, arg GetCardParams) (Card, error)
	GetCards(ctx context.Context, accountID int32) ([]Card, error)
	GetTenant(ctx context.Context, id int32) (Tenant, error)
	GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error)
	GetTransactions(ctx context.Context, cardID int32) ([]Transaction, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
