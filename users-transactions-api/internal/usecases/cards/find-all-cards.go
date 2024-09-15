package usecases

import (
	"context"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	usecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/account"
)

type FindAllCards struct {
	repo               infra.Querier
	findAccountUsecase *usecases.FindOneAccountUsecase
}

func NewFindAllCards(repo infra.Querier,
	findAccountUsecase *usecases.FindOneAccountUsecase) *FindAllCards {
	return &FindAllCards{
		repo:               repo,
		findAccountUsecase: findAccountUsecase,
	}
}

func (uc *FindAllCards) FindAll(tenantId int32, accountId int32) ([]infra.Card, error) {
	_, err := uc.findAccountUsecase.FindOne(tenantId, accountId)

	if err != nil {
		return nil, err
	}

	cards := make([]infra.Card, 0)

	result, err := uc.repo.GetCards(context.Background(), accountId)

	if err != nil {
		slog.Error(
			"error to find all cards",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	cards = append(cards, result...)

	return cards, nil
}
