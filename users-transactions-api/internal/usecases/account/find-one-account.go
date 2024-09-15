package usecases

import (
	"context"
	"database/sql"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
)

type FindOneAccountUsecase struct {
	repo infra.Querier
}

func NewFindOneAccountUsecase(repo infra.Querier) *FindOneAccountUsecase {
	return &FindOneAccountUsecase{
		repo: repo,
	}
}

func (uc *FindOneAccountUsecase) FindOne(tenantId int32, accountId int32) (*infra.Account, error) {
	account, err := uc.repo.GetAccount(context.Background(), infra.GetAccountParams{
		TenantID: tenantId,
		ID:       accountId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &shared.EntityNotFoundError{
				Object: "account",
				Id:     accountId,
			}
		}
		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return &account, nil
}
