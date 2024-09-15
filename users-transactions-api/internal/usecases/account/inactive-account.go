package usecases

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
)

type InactiveAccountUsecase struct {
	repo infra.Querier
}

func NewInactiveAccountUsecase(repo infra.Querier) *InactiveAccountUsecase {
	return &InactiveAccountUsecase{
		repo: repo,
	}
}

func (uc *InactiveAccountUsecase) Inactive(tenantId int32, id int32) error {
	ctx := context.Background()
	account, err := uc.repo.GetAccount(ctx, infra.GetAccountParams{
		TenantID: tenantId,
		ID:       id,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return &shared.EntityNotFoundError{
				Object: "account",
				Id:     id,
			}
		}
		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return err
	}

	currentTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	_, err = uc.repo.UpdateAccount(ctx, infra.UpdateAccountParams{
		ID:        account.ID,
		Status:    "inactive",
		UpdatedAt: currentTime,
		DeletedAt: currentTime,
	})

	if err != nil {
		return err
	}

	return nil
}
