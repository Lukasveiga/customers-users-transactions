package infra

import (
	"fmt"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgreCardRepository(t *testing.T) {
	dbClient := config.SetupPgTestcontainers()

	pgCardRepository := NewPgCardRepository(dbClient)
	pgAccountRepository := NewPgAccountRepository(dbClient)

	defer dbClient.Close()

	account := &domain.Account{
		TenantId: 1,
		Status:   "active",
	}

	card := &domain.Card{
		Amount: 0,
	}

	t.Run("should successfully connect to pg container", func(t *testing.T) {
		assert.NotNil(t, &dbClient)
	})

	t.Run("[Save] should save card and return it", func(t *testing.T) {
		savedAccount, err := pgAccountRepository.Save(account)

		assert.NoError(t, err)

		card.AccountId = savedAccount.Id

		savedCard, err := pgCardRepository.Save(card)

		assert.NoError(t, err)
		assert.Equal(t, savedAccount.Id, savedCard.AccountId)
		assert.Equal(t, card.Amount, savedCard.Amount)
	})

	t.Run("[Save] should return error when account id not exists", func(t *testing.T) {
		card.AccountId = 5

		savedCard, err := pgCardRepository.Save(card)

		assert.Nil(t, savedCard)
		assert.Contains(t, err.(*pq.Error).Detail,
			fmt.Sprintf("Key (account_id)=(%d) is not present in table \"accounts\".", 5))
	})

	t.Run("[FindById] should find card by id and account id", func(t *testing.T) {
		savedAccount, err := pgAccountRepository.Save(account)
		assert.NoError(t, err)

		card.AccountId = savedAccount.Id
		savedCard, err := pgCardRepository.Save(card)
		assert.NoError(t, err)

		foundCard, err := pgCardRepository.FindById(savedAccount.Id, savedCard.Id)
		assert.NoError(t, err)

		assert.Equal(t, savedAccount.Id, foundCard.AccountId)
		assert.Equal(t, card.Amount, foundCard.Amount)
	})

	t.Run("[FindById] should return nil when card is not found", func(t *testing.T) {
		savedAccount, err := pgAccountRepository.Save(account)
		assert.NoError(t, err)

		card.AccountId = savedAccount.Id
		savedCard, err := pgCardRepository.Save(card)
		assert.NoError(t, err)

		testCases := []struct {
			name      string
			accountId int32
			cardId    int32
		}{
			{
				name:      "invalid account id",
				accountId: 99,
				cardId:    savedCard.Id,
			},
			{
				name:      "invalid card id",
				accountId: savedAccount.Id,
				cardId:    99,
			},
		}

		for i := range testCases {
			tc := testCases[i]

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				foundCard, err := pgCardRepository.FindById(tc.accountId, tc.cardId)

				assert.NoError(t, err)
				assert.Nil(t, foundCard)
			})
		}
	})
}
