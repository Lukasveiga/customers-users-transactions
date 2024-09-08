package usecases

import (
	"context"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
)

type FindAllUsecase struct {
	repo infra.Querier
}

func NewFindAllAccountsUsecase(repo infra.Querier) *FindAllUsecase {
	return &FindAllUsecase{
		repo: repo,
	}
}

func (uc *FindAllUsecase) FindAll(tenantId int32) ([]infra.Account, error) {
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
