package usecases

import (
	"context"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
)

type CreateCardUsecase struct {
	repo               infra.Querier
	findAccountUsecase *usecases.FindOneAccountUsecase
}

func NewCreateCardUsecase(repo infra.Querier,
	findAccountUsecase *usecases.FindOneAccountUsecase) *CreateCardUsecase {
	return &CreateCardUsecase{
		repo:               repo,
		findAccountUsecase: findAccountUsecase,
	}
}

func (uc *CreateCardUsecase) Create(tenantId int32, accountId int32) (*infra.Card, error) {
	account, err := uc.findAccountUsecase.FindOne(tenantId, accountId)

	if err != nil {
		return nil, err
	}

	if account.Status == "inactive" {
		return nil, &shared.InactiveAccountError{}
	}

	savedCard, err := uc.repo.CreateCard(context.Background(), accountId)

	if err != nil {
		slog.Error(
			"error creating card",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return &savedCard, nil
}
