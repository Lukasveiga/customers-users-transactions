package usecases

import (
	"github.com/Lukasveiga/customers-users-Transaction/config"
	"github.com/Lukasveiga/customers-users-Transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-Transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-Transaction/internal/shared"
)

type CreateAccountUsecase struct {
	repo port.AccountRepository
	logger config.Logger
}

func NewCreateAccountUsecase(repo port.AccountRepository, logger config.Logger) *CreateAccountUsecase {
	return &CreateAccountUsecase{
		repo:  repo,
		logger: logger,
	}
}

func (uc CreateAccountUsecase) Exec(account *domain.Account) (*domain.Account, error) {
	existAccount, err := uc.repo.FindByNumber(account.TenantId, account.Number);

	if err != nil {
		uc.logger.Errorf("error to find account by id: %v", err)
		return nil, err
	}

	if existAccount != nil {
		return nil, &shared.AlreadyExistsError{
			Object: "account",
			Id: account.Number,
		}
	}

	savedAccount, err := uc.repo.Create(account)

	if err != nil {
		uc.logger.Errorf("error creating account: %v", err)
		return nil, err
	}
	
	return savedAccount, nil
} 