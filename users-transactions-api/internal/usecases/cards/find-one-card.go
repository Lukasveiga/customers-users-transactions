package usecases

import (
	"context"
	"database/sql"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
)

type FindCardUsecase struct {
	repo               infra.Querier
	findAccountUsecase *usecases.FindOneAccountUsecase
}

func NewFindCardUsecase(repo infra.Querier,
	findAccountUsecase *usecases.FindOneAccountUsecase) *FindCardUsecase {
	return &FindCardUsecase{
		repo:               repo,
		findAccountUsecase: findAccountUsecase,
	}
}

func (uc *FindCardUsecase) FindOne(tenantId int32, accountId int32, cardId int32) (*infra.Card, error) {
	_, err := uc.findAccountUsecase.FindOne(tenantId, accountId)

	if err != nil {
		return nil, err
	}

	card, err := uc.repo.GetCard(context.Background(), infra.GetCardParams{
		AccountID: accountId,
		ID:        cardId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &shared.EntityNotFoundError{
				Object: "card",
				Id:     cardId,
			}
		}
		slog.Error(
			"error creating card",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return &card, nil
}
