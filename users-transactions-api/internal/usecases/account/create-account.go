package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
)

type CreateAccountUsecase struct {
	repo port.AccountRepository
}

func NewCreateAccountUsecase(repo port.AccountRepository) *CreateAccountUsecase {
	return &CreateAccountUsecase{
		repo: repo,
	}
}

func (uc *CreateAccountUsecase) Create(tenantId int32) (*domain.Account, error) {
	account := &domain.Account{
		TenantId: tenantId,
	}

	savedAccount, err := uc.repo.Save(account.Create())

	if err != nil {
		slog.Error(
			"error creating account",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return savedAccount, nil
}
