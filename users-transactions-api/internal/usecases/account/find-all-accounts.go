package usecases

import (
	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
)

type FindAllUsecase struct {
	repo port.AccountRepository
	logger config.Logger
}

func NewFindAllAccountsUsecase(repo port.AccountRepository, logger config.Logger) *FindAllUsecase {
	return &FindAllUsecase{
		repo: repo,
		logger: logger,
	}
}

func (uc FindAllUsecase) Exec(tenantId int32) ([]domain.Account, error) {
	accounts, err := uc.repo.FindAll(tenantId)

	if err != nil {
		uc.logger.Errorf("error to find all accounts: %v", err)
		return nil, err
	}

	return accounts, nil;
}