package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/mocks"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/account"
	cardUsecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/cards"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCardHandler(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)

	findOneAccountUsecase := accountUsecases.NewFindOneAccountUsecase(mockRepo)
	createCardUsecase := cardUsecases.NewCreateCardUsecase(mockRepo, findOneAccountUsecase)
	findCardUsecase := cardUsecases.NewFindCardUsecase(mockRepo, findOneAccountUsecase)
	findAllCardsUsecase := cardUsecases.NewFindAllCards(mockRepo, findOneAccountUsecase)

	sut := NewCardHandler(createCardUsecase, findCardUsecase, findAllCardsUsecase)

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

	t.Run("[Create] invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[Create] invalid account id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "invalid",
		}}

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid account id", responseBody["error"])
	})

	t.Run("[Create] Internal server error", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[Create] Error account not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		accountId := "1"

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: accountId,
		}}

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %s", accountId),
			responseBody["error"])
	})

	t.Run("[Create] Error account inactive", func(t *testing.T) {
		account.Status = "inactive"

		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		}}

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Card cannot be created to an inactive account",
			responseBody["error"])
	})

	t.Run("[Create] Card created successfully", func(t *testing.T) {
		account.Status = "active"

		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("CreateCard").Return(card, nil)
		defer mockRepo.On("CreateCard").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		}}

		sut.Create(c)

		var responseBody infra.Card
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, card, responseBody)
	})

	t.Run("[FindCard] Error invalid card id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		},
			{
				Key:   "cardId",
				Value: "invalid",
			}}

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid card id", responseBody["error"])
	})

	t.Run("[FindCard] Error account not found", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		accountId := "1"

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: accountId,
		},
			{
				Key:   "cardId",
				Value: "1",
			}}

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %s", accountId),
			responseBody["error"])
	})

	t.Run("[FindCard] Success to find card by id", func(t *testing.T) {
		account.Status = "active"

		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("GetCard").Return(card, nil)
		defer mockRepo.On("GetCard").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		},
			{
				Key:   "cardId",
				Value: fmt.Sprint(card.ID),
			}}

		sut.FindOne(c)

		var responseBody infra.Card
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, card, responseBody)
	})

	t.Run("[FindAll] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/cards", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[FindAll] Invalid account id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/cards", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "invalid",
		}}

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid account id", responseBody["error"])
	})

	t.Run("[FindAll] Internal server error", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, errors.New("internal error"))
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/cards", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		}}

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[FindAll] Account not found error", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/cards", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		}}

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %d", account.ID),
			responseBody["error"])
	})

	t.Run("[FindAll] Success", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		cards := []infra.Card{
			card,
		}

		mockRepo.On("GetCards").Return(cards, nil)
		defer mockRepo.On("GetCards").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/cards", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.ID),
		}}

		sut.FindAll(c)

		var responseBody []infra.Card
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Len(t, responseBody, 1)
		assert.Equal(t, cards, responseBody)
	})

}
