package usecases

import (
	"log/slog"

	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type InactiveAccountUsecase struct {
	repo port.AccountRepository
}

func NewInactiveAccountUsecase(repo port.AccountRepository) *InactiveAccountUsecase {
	return &InactiveAccountUsecase{
		repo: repo,
	}
}

func (uc *InactiveAccountUsecase) Inactive(tenantId int32, id int32) error {
	account, err := uc.repo.FindById(tenantId, id)

	if err != nil {
		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return err
	}

	if account == nil {
		return &shared.EntityNotFoundError{
			Object: "account",
			Id:     id,
		}
	}

	account.Inactive()

	_, err = uc.repo.Update(account)

	if err != nil {
		return err
	}

	return nil
}
