package usecases

import (
	"context"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
)

type CreateAccountUsecase struct {
	repo infra.Querier
}

func NewCreateAccountUsecase(repo infra.Querier) *CreateAccountUsecase {
	return &CreateAccountUsecase{
		repo: repo,
	}
}

func (uc *CreateAccountUsecase) Create(tenantId int32) (*infra.Account, error) {
	savedAccount, err := uc.repo.CreateAccount(context.Background(), infra.CreateAccountParams{
		TenantID: tenantId,
		Status:   "active",
	})

	if err != nil {
		slog.Error(
			"error creating account",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return &savedAccount, nil
}
