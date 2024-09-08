package usecases

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type ActiveAccountUsecase struct {
	repo infra.Querier
}

func NewActiveAccountUsecase(repo infra.Querier) *ActiveAccountUsecase {
	return &ActiveAccountUsecase{
		repo: repo,
	}
}

func (uc *ActiveAccountUsecase) Active(tenantId int32, accountId int32) error {
	ctx := context.Background()
	account, err := uc.repo.GetAccount(ctx, infra.GetAccountParams{
		TenantID: tenantId,
		ID:       accountId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return &shared.EntityNotFoundError{
				Object: "account",
				Id:     accountId,
			}
		}
		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return err
	}

	_, err = uc.repo.UpdateAccount(ctx, infra.UpdateAccountParams{
		ID:     account.ID,
		Status: "active",
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		DeletedAt: account.DeletedAt,
	})

	if err != nil {
		slog.Error(
			"error when update account",
			slog.String("err", err.Error()),
		)
		return err
	}

	return nil
}
