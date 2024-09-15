package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/handlers/dto"
	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/mocks"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/account"
	cardUsecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/cards"
	transactionUsecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/transactions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHandler(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)

	findAccountUsecase := accountUsecases.NewFindOneAccountUsecase(mockRepo)
	findCardUsecase := cardUsecases.NewFindCardUsecase(mockRepo, findAccountUsecase)
	createTransactionUsecase := transactionUsecases.NewCreateTransactionUsecase(mockRepo, findCardUsecase)
	findTransactionUsecase := transactionUsecases.NewFindTransactionUsecase(mockRepo, findCardUsecase)
	findTransactionsUsecase := transactionUsecases.NewFindAllTransactionsUsecase(mockRepo, findCardUsecase)

	sut := NewTransactionHandler(createTransactionUsecase, findTransactionUsecase, findTransactionsUsecase)

	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	card := infra.Card{
		ID:        1,
		Amount:    0,
		AccountID: 1,
	}

	transaction := infra.Transaction{
		ID:     1,
		CardID: 1,
		Kind:   "Streaming Z",
		Value:  50,
	}

	t.Run("[Create] Transaction created successfully", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("CreateTransactionTx").Return(transaction, nil)
		defer mockRepo.On("CreateTransactionTx").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		body, err := json.Marshal(dto.TransactioRequest{
			CardId: transaction.CardID,
			Kind:   transaction.Kind,
			Value:  transaction.Value,
		})

		assert.NoError(t, err)

		c.Request = httptest.NewRequest("POST", "/transaction", bytes.NewReader(body))
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		}}

		sut.Create(c)

		var responseBody dto.TransactionResponse
		err = json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, dto.TransactionResponse{
			ID:     transaction.ID,
			CardId: transaction.CardID,
			Kind:   transaction.Kind,
			Value:  transaction.Value,
		}, responseBody)
	})

	t.Run("[FindOne] Success to find a transaction", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("GetTransaction").Return(transaction, nil)
		defer mockRepo.On("GetTransaction").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transaction", nil)
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{
			{
				Key:   "accountId",
				Value: fmt.Sprint(account.ID),
			},
			{
				Key:   "cardId",
				Value: fmt.Sprint(card.ID),
			},
			{
				Key:   "transactionId",
				Value: fmt.Sprint(transaction.ID),
			},
		}

		sut.FindOne(c)

		var responseBody dto.TransactionResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, dto.TransactionResponse{
			ID:     transaction.ID,
			CardId: transaction.CardID,
			Kind:   transaction.Kind,
			Value:  transaction.Value,
		}, responseBody)
	})

	t.Run("[FindOne] Error transaction not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("GetTransaction").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetTransaction").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transaction", nil)
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{
			{
				Key:   "accountId",
				Value: fmt.Sprint(account.ID),
			},
			{
				Key:   "cardId",
				Value: fmt.Sprint(card.ID),
			},
			{
				Key:   "transactionId",
				Value: fmt.Sprint(transaction.ID),
			},
		}

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("transaction not found with id %d", transaction.ID),
			responseBody["error"])
	})

	t.Run("[FindAll] Success to find all transactions by card id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		mockRepo.On("GetTransactions").Return([]infra.Transaction{
			transaction,
		}, nil)
		defer mockRepo.On("GetTransactions").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transaction", nil)
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{
			{
				Key:   "accountId",
				Value: fmt.Sprint(account.ID),
			},
			{
				Key:   "cardId",
				Value: fmt.Sprint(card.ID),
			},
		}

		sut.FindAll(c)

		var responseBody []dto.TransactionResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, []dto.TransactionResponse{
			{
				ID:     transaction.ID,
				CardId: transaction.CardID,
				Kind:   transaction.Kind,
				Value:  transaction.Value,
			},
		}, responseBody)
	})

	t.Run("[FindAll] Error card not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetCard").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transaction", nil)
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{
			{
				Key:   "accountId",
				Value: fmt.Sprint(account.ID),
			},
			{
				Key:   "cardId",
				Value: fmt.Sprint(card.ID),
			},
		}

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("card not found with id %d", transaction.ID),
			responseBody["error"])
	})
}
