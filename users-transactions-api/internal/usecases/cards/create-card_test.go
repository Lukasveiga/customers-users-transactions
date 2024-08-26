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

func TestCreateCardUsecase(t *testing.T) {
	t.Parallel()

	mockCardRepo := new(mocks.MockCardRepository)
	mockAccountRepo := new(mocks.MockAccountRepository)

	findAccountUsecase := usecases.NewFindOneAccountUsecase(mockAccountRepo)

	sut := NewCreateCardUsecase(mockCardRepo, findAccountUsecase)

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

		result, err := sut.Create(1, card.AccountId)

		assert.Nil(t, result)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", card.AccountId),
			err.Error())
	})

	t.Run("Error to save card", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("Save").Return(nil, errors.New("Internal error"))
		defer mockCardRepo.On("Save").Unset()

		result, err := sut.Create(1, card.AccountId)

		assert.Nil(t, result)
		assert.Equal(t, "Internal error", err.Error())
	})

	t.Run("Success to create card", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("Save").Return(card, nil)
		defer mockCardRepo.On("Save").Unset()

		result, err := sut.Create(1, card.AccountId)

		assert.NoError(t, err)
		assert.Equal(t, card, result)
		assert.Equal(t, card.Amount, result.Amount)
	})

	t.Run("Inactive account error", func(t *testing.T) {
		account.Status = "inactive"

		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		result, err := sut.Create(1, card.AccountId)

		assert.Nil(t, result)
		assert.Equal(t, "Card cannot be created to an inactive account", err.Error())
	})
}
