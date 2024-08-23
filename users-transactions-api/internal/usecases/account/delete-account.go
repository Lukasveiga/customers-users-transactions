package usecases

import (
	"log/slog"

	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type DeleteAccountUsecase struct {
	repo port.AccountRepository
}

func NewDeleteAccountUsecase(repo port.AccountRepository) *DeleteAccountUsecase {
	return &DeleteAccountUsecase{
		repo: repo,
	}
}

func (uc *DeleteAccountUsecase) Delete(tenantId int32, id int32) error {
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

	err = uc.repo.Delete(tenantId, id)

	if err != nil {
		return err
	}

	return nil
}
