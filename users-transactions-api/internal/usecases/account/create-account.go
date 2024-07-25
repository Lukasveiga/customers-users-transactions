package usecases

import (
	"log/slog"
	"time"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type CreateAccountUsecase struct {
	repo port.AccountRepository
}

func NewCreateAccountUsecase(repo port.AccountRepository) *CreateAccountUsecase {
	return &CreateAccountUsecase{
		repo: repo,
	}
}

func (uc CreateAccountUsecase) Exec(account *domain.Account) (*domain.Account, error) {
	existAccount, err := uc.repo.FindByNumber(account.TenantId, account.Number)

	if err != nil {
		slog.Error(
			"error to find account by id",
			slog.Int("id", int(account.Id)),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if existAccount != nil {
		return nil, &shared.AlreadyExistsError{
			Object: "account",
			Id:     account.Number,
		}
	}

	err = account.Validate()
	if err != nil {
		return nil, err
	}

	account.CreatedAt = time.Now().UTC()
	savedAccount, err := uc.repo.Create(account)

	if err != nil {
		slog.Error(
			"error creating account",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return savedAccount, nil
}
