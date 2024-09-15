package usecases

import (
	"context"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/cards"
)

type CreateTransactionUsecase struct {
	repo            infra.QuerierTx
	findCardUsecase *usecases.FindCardUsecase
}

func NewCreateTransactionUsecase(repo infra.QuerierTx,
	findCardUsecase *usecases.FindCardUsecase) *CreateTransactionUsecase {
	return &CreateTransactionUsecase{
		repo:            repo,
		findCardUsecase: findCardUsecase,
	}
}

func (uc *CreateTransactionUsecase) Create(tenantId int32, accountId int32,
	transaction infra.Transaction) (*infra.Transaction, error) {
	card, err := uc.findCardUsecase.FindOne(tenantId, accountId, transaction.CardID)

	if err != nil {
		return nil, err
	}

	err = transactionInputValidation(transaction)

	if err != nil {
		return nil, err
	}

	savedTransaction, err := uc.repo.CreateTransactionTx(context.Background(), infra.CreateTransactionParams{
		CardID: card.ID,
		Kind:   transaction.Kind,
		Value:  transaction.Value,
	})

	if err != nil {
		return nil, err
	}

	return &savedTransaction, nil
}

func transactionInputValidation(t infra.Transaction) error {
	valErr := &shared.ValidationError{
		Errors: make(map[string]string),
	}

	if len(t.Kind) == 0 {
		valErr.AddError("kind", "cannot be empty")
	}

	if t.Value < 0 {
		valErr.AddError("value", "must be greater than zero (0)")
	}

	if valErr.HasErrors() {
		return valErr
	}

	return nil
}
