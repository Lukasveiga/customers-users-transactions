package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	cardUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
	"github.com/stretchr/testify/assert"
)

func TestFindTransactionUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	fincAccountUsecase := accountUsecases.NewFindOneAccountUsecase(mockRepo)
	findCardUsecase := cardUsecases.NewFindCardUsecase(mockRepo, fincAccountUsecase)

	sut := NewFindTransactionUsecase(mockRepo, findCardUsecase)

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

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("GetTransaction").Return(transaction, nil)
		defer mockRepo.On("GetTransaction").Unset()

		result, err := sut.FindOne(1, 1, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, transaction, *result)
	})

	t.Run("Error account not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.FindOne(1, 1, 1, 1)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", transaction.CardID), err.Error())
	})

	t.Run("Error card not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetCard").Unset()

		result, err := sut.FindOne(1, 1, 1, 1)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("card not found with id %d", transaction.CardID), err.Error())
	})

	t.Run("Error transaction not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("GetTransaction").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetTransaction").Unset()

		result, err := sut.FindOne(1, 1, 1, 1)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("transaction not found with id %d", transaction.CardID), err.Error())
	})

	t.Run("Internal error", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("GetTransaction").Return(nil, errors.New("internal error"))
		defer mockRepo.On("GetTransaction").Unset()

		result, err := sut.FindOne(1, 1, 1, 1)

		assert.Nil(t, result)
		assert.EqualError(t, errors.New("internal error"), err.Error())
	})
}
