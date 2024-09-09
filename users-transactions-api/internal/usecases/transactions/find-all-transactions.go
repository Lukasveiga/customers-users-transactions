package usecases

import (
	"context"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
)

type FindAllTransactionsUsecase struct {
	repo            infra.Querier
	findCardUsecase *usecases.FindCardUsecase
}

func NewFindAllTransactionsUsecase(repo infra.Querier,
	findCardUsecase *usecases.FindCardUsecase) *FindAllTransactionsUsecase {
	return &FindAllTransactionsUsecase{
		repo:            repo,
		findCardUsecase: findCardUsecase,
	}
}

func (uc *FindAllTransactionsUsecase) FindAll(tenantId int32, accountId int32,
	cardId int32) ([]infra.Transaction, error) {
	_, err := uc.findCardUsecase.FindOne(tenantId, accountId, cardId)

	if err != nil {
		return nil, err
	}

	transactions := make([]infra.Transaction, 0)

	result, err := uc.repo.GetTransactions(context.Background(), cardId)

	if err != nil {
		slog.Error(
			"error to find all transactions",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	transactions = append(transactions, result...)

	return transactions, nil
}
