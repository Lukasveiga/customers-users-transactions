package usecases

import (
	"context"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
)

type FindAllAccountsUsecase struct {
	repo infra.Querier
}

func NewFindAllAccountsUsecase(repo infra.Querier) *FindAllAccountsUsecase {
	return &FindAllAccountsUsecase{
		repo: repo,
	}
}

func (uc *FindAllAccountsUsecase) FindAll(tenantId int32) ([]infra.Account, error) {
	accounts := make([]infra.Account, 0)

	result, err := uc.repo.GetAccounts(context.Background(), tenantId)

	if err != nil {
		slog.Error(
			"error to find all accounts",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	accounts = append(accounts, result...)

	return accounts, nil
}
