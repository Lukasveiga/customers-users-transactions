package usecases

import (
	"testing"

	"github.com/Lukasveiga/customers-users-Transaction/config"
	"github.com/Lukasveiga/customers-users-Transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-Transaction/internal/mocks"
	"github.com/Lukasveiga/customers-users-Transaction/internal/shared"
	"github.com/stretchr/testify/assert"
)



func TestCreateAccountUsecase(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	sut := NewCreateAccountUsecase(mockRepo, *config.NewLogger("test"))

	t.Run("Error account number already exists", func(t *testing.T) {
		account := &domain.Account{
			Number: "123",
		}

		mockRepo.On("FindByNumber").Return(account, nil)

		expectedErr := &shared.AlreadyExistsError{
			Object: "account",
			Id: account.Number,
		}

		result, err := sut.Exec(account)

		assert.Nil(t, result)
		assert.Equal(t, err.Error(), expectedErr.Error())
	})
}