package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
)

type CreateCardUsecase struct {
	repo               port.CardRepository
	findAccountUsecase *usecases.FindOneAccountUsecase
}

func NewCreateCardUsecase(repo port.CardRepository,
	findAccountUsecase *usecases.FindOneAccountUsecase) *CreateCardUsecase {
	return &CreateCardUsecase{
		repo:               repo,
		findAccountUsecase: findAccountUsecase,
	}
}

func (uc *CreateCardUsecase) Create(tenantId int32, accountId int32) (*domain.Card, error) {
	account, err := uc.findAccountUsecase.FindOne(tenantId, accountId)

	if err != nil {
		return nil, err
	}

	if account.Status == "inactive" {
		return nil, &shared.InactiveAccountError{}
	}

	card := &domain.Card{
		AccountId: accountId,
	}

	savedCard, err := uc.repo.Save(card.Create())

	if err != nil {
		slog.Error(
			"error creating card",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return savedCard, nil
}
