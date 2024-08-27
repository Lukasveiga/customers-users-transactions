package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
)

type FindCardUsecase struct {
	repo               port.CardRepository
	findAccountUsecase *usecases.FindOneAccountUsecase
}

func NewFindCardUsecase(repo port.CardRepository,
	findAccountUsecase *usecases.FindOneAccountUsecase) *FindCardUsecase {
	return &FindCardUsecase{
		repo:               repo,
		findAccountUsecase: findAccountUsecase,
	}
}

func (uc *FindCardUsecase) FindOne(tenantId int32, accountId int32, cardId int32) (*domain.Card, error) {
	_, err := uc.findAccountUsecase.FindOne(tenantId, accountId)

	if err != nil {
		return nil, err
	}

	card, err := uc.repo.FindById(accountId, cardId)

	if err != nil {
		slog.Error(
			"error creating card",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if card == nil {
		return nil, &shared.EntityNotFoundError{
			Object: "Card",
			Id:     cardId,
		}
	}

	return card, nil
}
