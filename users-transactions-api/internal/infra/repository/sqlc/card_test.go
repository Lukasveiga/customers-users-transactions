package infra

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestCard(t *testing.T, accountId int32) Card {
	arg := accountId

	card, err := testQueries.CreateCard(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, card)

	return card
}

func TestCardRepository(t *testing.T) {

	t.Run("[CreateCard] should create new card and return it", func(t *testing.T) {
		account := createTestAccount(t, 1)
		createTestCard(t, account.ID)
	})

	t.Run("[GetCard] should find card by id and account id", func(t *testing.T) {
		account := createTestAccount(t, 1)
		card := createTestCard(t, account.ID)

		foundCard, err := testQueries.GetCard(context.Background(), GetCardParams{
			AccountID: account.ID,
			ID:        card.ID,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, foundCard)
		assert.Equal(t, card, foundCard)
	})

	t.Run("[GetCards] should return a list of card given a account id", func(t *testing.T) {
		n := 5
		account := createTestAccount(t, 2)

		for i := 0; i < n; i++ {
			createTestCard(t, account.ID)
		}

		cards, err := testQueries.GetCards(context.Background(), account.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, cards)
		assert.Len(t, cards, n)

		for _, card := range cards {
			assert.NotEmpty(t, card)
		}
	})

	t.Run("[AddAmount] should update card and return it", func(t *testing.T) {
		account := createTestAccount(t, 3)
		card := createTestCard(t, account.ID)

		arg := AddAmountParams{
			Amount: 200,
			ID:     card.ID,
			UpdatedAt: sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: true,
			},
		}

		updatedCard, err := testQueries.AddAmount(context.Background(), arg)

		assert.NoError(t, err)
		assert.NotEmpty(t, updatedCard)
		assert.Equal(t, arg.Amount, updatedCard.Amount)
	})
}
