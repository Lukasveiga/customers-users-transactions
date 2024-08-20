package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
)

type FindAllUsecase struct {
	repo port.AccountRepository
}

func NewFindAllAccountsUsecase(repo port.AccountRepository) *FindAllUsecase {
	return &FindAllUsecase{
		repo: repo,
	}
}

func (uc *FindAllUsecase) FindAll(tenantId int32) ([]*domain.Account, error) {
	accounts := make([]*domain.Account, 0)

	result, err := uc.repo.FindAll(tenantId)

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
