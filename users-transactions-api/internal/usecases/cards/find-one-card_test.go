package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/account"
	"github.com/stretchr/testify/assert"
)

func TestFindCardUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)

	findAccountUsecase := usecases.NewFindOneAccountUsecase(mockRepo)

	sut := NewFindCardUsecase(mockRepo, findAccountUsecase)

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

	t.Run("Error to find account", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.FindOne(1, 1, card.AccountID)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", card.AccountID),
			err.Error())
	})

	t.Run("Error to find card by id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(nil, errors.New("Internal error"))
		defer mockRepo.On("GetCard").Unset()

		result, err := sut.FindOne(1, 1, card.AccountID)

		assert.Nil(t, result)
		assert.Equal(t, "Internal error", err.Error())
	})

	t.Run("Erro card not found by id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetCard").Unset()

		result, err := sut.FindOne(1, 1, card.AccountID)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("card not found with id %d", card.ID),
			err.Error())
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		result, err := sut.FindOne(1, 1, card.AccountID)

		assert.NoError(t, err)
		assert.Equal(t, &card, result)
	})
}
