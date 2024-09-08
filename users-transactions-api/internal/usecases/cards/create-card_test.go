package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	"github.com/stretchr/testify/assert"
)

func TestCreateCardUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)

	findAccountUsecase := usecases.NewFindOneAccountUsecase(mockRepo)

	sut := NewCreateCardUsecase(mockRepo, findAccountUsecase)

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

		result, err := sut.Create(1, card.AccountID)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", card.AccountID),
			err.Error())
	})

	t.Run("Error to save card", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("CreateCard").Return(nil, errors.New("Internal error"))
		defer mockRepo.On("CreateCard").Unset()

		result, err := sut.Create(1, card.AccountID)

		assert.Nil(t, result)
		assert.Equal(t, "Internal error", err.Error())
	})

	t.Run("Success to create card", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("CreateCard").Return(card, nil)
		defer mockRepo.On("CreateCard").Unset()

		result, err := sut.Create(1, card.AccountID)

		assert.NoError(t, err)
		assert.Equal(t, &card, result)
		assert.Equal(t, card.Amount, result.Amount)
	})

	t.Run("Inactive account error", func(t *testing.T) {
		account.Status = "inactive"

		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.Create(1, card.AccountID)

		assert.Nil(t, result)
		assert.Equal(t, "Card cannot be created to an inactive account", err.Error())
	})
}
