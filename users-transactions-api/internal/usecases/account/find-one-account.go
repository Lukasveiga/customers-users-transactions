package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type FindOneAccountUsecase struct {
	repo port.AccountRepository
}

func NewFindOneAccountUsecase(repo port.AccountRepository) *FindOneAccountUsecase {
	return &FindOneAccountUsecase{
		repo: repo,
	}
}

func (uc *FindOneAccountUsecase) FindOne(tenantId int32, accountId int32) (*domain.Account, error) {
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

	return account, nil
}
