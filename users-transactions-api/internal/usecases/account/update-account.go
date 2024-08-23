package usecases

import (
	"log/slog"
	"time"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type UpdateAccountUsecase struct {
	repo port.AccountRepository
}

func NewUpdateAccountUsecase(repo port.AccountRepository) *UpdateAccountUsecase {
	return &UpdateAccountUsecase{
		repo: repo,
	}
}

func (uc *UpdateAccountUsecase) Update(accountId int32, account *domain.Account) (*domain.Account, error) {
	existingAccount, err := uc.repo.FindById(account.TenantId, accountId)

	if err != nil {
		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if existingAccount == nil {
		return nil, &shared.EntityNotFoundError{
			Object: "account",
			Id:     accountId,
		}
	}

	existingAccountNumber, err := uc.repo.FindByNumber(account.TenantId, account.Number)

	if err != nil {
		slog.Error(
			"error to find account by number",
			slog.String("number", account.Number),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if existingAccountNumber != nil && existingAccountNumber.Id != accountId && account.Number == existingAccountNumber.Number {
		return nil, &shared.AlreadyExistsError{
			Object: "account",
			Id:     account.Number,
		}
	}

	err = account.Validate()
	if err != nil {
		return nil, err
	}

	existingAccount.UpdatedAt = time.Now().UTC()
	existingAccount.Number = account.Number
	existingAccount.Status = account.Status

	updatedAccount, err := uc.repo.Update(accountId, existingAccount)

	if err != nil {
		slog.Error(
			"error when update account",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return updatedAccount, nil
}
