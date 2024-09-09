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

func TestFindAllCards(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)

	findAccountUsecase := usecases.NewFindOneAccountUsecase(mockRepo)

	sut := NewFindAllCards(mockRepo, findAccountUsecase)

	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	cards := []infra.Card{
		{
			ID:        1,
			Amount:    200,
			AccountID: 1,
		},
		{
			ID:        1,
			Amount:    200,
			AccountID: 1,
		},
	}

	t.Run("Error account not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		result, err := sut.FindAll(1, account.ID)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", account.ID),
			err.Error())
	})

	t.Run("Error to get cards", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCards").Return(nil, errors.New("internal error"))
		defer mockRepo.On("GetCards").Unset()

		result, err := sut.FindAll(1, account.ID)

		assert.Nil(t, result)
		assert.EqualError(t, errors.New("internal error"), err.Error())
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCards").Return(cards, nil)
		defer mockRepo.On("GetCards").Unset()

		result, err := sut.FindAll(1, account.ID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		for _, card := range result {
			assert.NotEmpty(t, card)
		}
	})
}
