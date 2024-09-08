package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	cardUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	fincAccountUsecase := accountUsecases.NewFindOneAccountUsecase(mockRepo)
	findCardUsecase := cardUsecases.NewFindCardUsecase(mockRepo, fincAccountUsecase)

	sut := NewCreateTransactionUsecase(mockRepo, findCardUsecase)

	card := infra.Card{
		ID:        1,
		Amount:    200,
		AccountID: 1,
	}

	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	transaction := infra.Transaction{
		ID:     1,
		CardID: 1,
		Kind:   "Streaming Z",
		Value:  200,
	}

	t.Run("Success to create transaction", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("CreateTransactionTx").Return(transaction, nil)
		defer mockRepo.On("CreateTransactionTx").Unset()

		savedTransaction, err := sut.Create(1, account.ID, transaction)

		assert.NoError(t, err)
		assert.Equal(t, &transaction, savedTransaction)
	})

	t.Run("Erro account not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		savedTransaction, err := sut.Create(1, account.ID, transaction)

		assert.Nil(t, savedTransaction)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", transaction.CardID), err.Error())
	})

	t.Run("Erro card not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetCard").Unset()

		savedTransaction, err := sut.Create(1, account.ID, transaction)

		assert.Nil(t, savedTransaction)
		assert.Equal(t, fmt.Sprintf("card not found with id %d", transaction.CardID), err.Error())
	})

	t.Run("Error input validation", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		invalidTransaction := infra.Transaction{
			ID:     1,
			CardID: 1,
			Kind:   "",
			Value:  -200,
		}

		expectedError := &shared.ValidationError{
			Errors: map[string]string{
				"kind":  "cannot be empty",
				"value": "must be greater than zero (0)",
			},
		}

		savedTransaction, err := sut.Create(1, account.ID, invalidTransaction)

		assert.Nil(t, savedTransaction)
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("Erro to create transaction", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("CreateTransactionTx").Return(nil, errors.New("internal error"))
		defer mockRepo.On("CreateTransactionTx").Unset()

		savedTransaction, err := sut.Create(1, account.ID, transaction)

		assert.Nil(t, savedTransaction)
		assert.EqualError(t, errors.New("internal error"), err.Error())
	})
}
