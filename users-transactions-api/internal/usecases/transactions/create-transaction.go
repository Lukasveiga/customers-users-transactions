package usecases

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
)

type CreateTransactionUsecase struct {
	repo            port.TransactionRepository
	findCardUsecase *usecases.FindCardUsecase
}

func NewCreateTransactionUsecase(repo port.TransactionRepository,
	findCardUsecase *usecases.FindCardUsecase) *CreateTransactionUsecase {
	return &CreateTransactionUsecase{
		repo:            repo,
		findCardUsecase: findCardUsecase,
	}
}

func (uc *CreateTransactionUsecase) Create(tenantId int32, accountId int32,
	transaction *domain.Transaction) (*domain.Transaction, error) {
	card, err := uc.findCardUsecase.FindOne(tenantId, accountId, transaction.CardId)

	if err != nil {
		return nil, err
	}

	validTransaction, err := transaction.Create()

	if err != nil {
		return nil, err
	}

	err = card.AddAmount(validTransaction.Value)

	if err != nil {
		return nil, err
	}

	// update card
	// save transaction

	return nil, nil
}
