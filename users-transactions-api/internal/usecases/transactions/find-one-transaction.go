package usecases

import (
	"context"
	"database/sql"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
)

type FindTransactionUsecase struct {
	repo            infra.Querier
	findCardUsecase *usecases.FindCardUsecase
}

func NewFindTransactionUsecase(repo infra.Querier, findCardUsecase *usecases.FindCardUsecase) *FindTransactionUsecase {
	return &FindTransactionUsecase{
		repo:            repo,
		findCardUsecase: findCardUsecase,
	}
}

func (uc *FindTransactionUsecase) FindOne(tenantId int32, accountId int32,
	cardId int32, transactionId int32) (*infra.Transaction, error) {
	_, err := uc.findCardUsecase.FindOne(tenantId, accountId, cardId)

	if err != nil {
		return nil, err
	}

	transaction, err := uc.repo.GetTransaction(context.Background(), infra.GetTransactionParams{
		CardID: cardId,
		ID:     transactionId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &shared.EntityNotFoundError{
				Object: "transaction",
				Id:     transactionId,
			}
		}
		slog.Error(
			"error to find transaction",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return &transaction, nil
}
