package usecases

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	"github.com/stretchr/testify/assert"
)

func TestFindCardUsecase(t *testing.T) {
	t.Parallel()

	mockCardRepo := new(mocks.MockCardRepository)
	mockAccountRepo := new(mocks.MockAccountRepository)

	findAccountUsecase := usecases.NewFindOneAccountUsecase(mockAccountRepo)

	sut := NewFindCardUsecase(mockCardRepo, findAccountUsecase)

	card := &domain.Card{
		Id:        1,
		Amount:    200,
		AccountId: 1,
	}

	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Status:   "active",
	}

	t.Run("Error to find account", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(nil, nil)
		defer mockAccountRepo.On("FindById").Unset()

		result, err := sut.FindOne(1, 1, card.AccountId)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", card.AccountId),
			err.Error())
	})

	t.Run("Error to find card by id", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("FindById").Return(nil, errors.New("Internal error"))
		defer mockCardRepo.On("FindById").Unset()

		result, err := sut.FindOne(1, 1, card.AccountId)

		assert.Nil(t, result)
		assert.Equal(t, "Internal error", err.Error())
	})

	t.Run("Erro card not found by id", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("FindById").Return(nil, nil)
		defer mockCardRepo.On("FindById").Unset()

		result, err := sut.FindOne(1, 1, card.AccountId)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("card not found with id %d", card.Id),
			err.Error())
	})

	t.Run("Success", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("FindById").Return(card, nil)
		defer mockCardRepo.On("FindById").Unset()

		result, err := sut.FindOne(1, 1, card.AccountId)

		assert.NoError(t, err)
		assert.Equal(t, card, result)
	})
}
