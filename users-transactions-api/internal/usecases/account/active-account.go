package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type ActiveAccountUsecase struct {
	repo port.AccountRepository
}

func NewActiveAccountUsecase(repo port.AccountRepository) *ActiveAccountUsecase {
	return &ActiveAccountUsecase{
		repo: repo,
	}
}

func (uc *ActiveAccountUsecase) Active(tenantId int32, accountId int32) (*domain.Account, error) {
	account, err := uc.repo.FindById(tenantId, accountId)

	if err != nil {
		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if account == nil {
		return nil, &shared.EntityNotFoundError{
			Object: "account",
			Id:     accountId,
		}
	}

	account.Active()

	updatedAccount, err := uc.repo.Update(account)

	if err != nil {
		slog.Error(
			"error when update account",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return updatedAccount, nil
}
