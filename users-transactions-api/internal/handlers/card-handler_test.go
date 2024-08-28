package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	cardUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCardHandler(t *testing.T) {
	t.Parallel()

	mockCardRepo := new(mocks.MockCardRepository)
	mockAccountRepo := new(mocks.MockAccountRepository)

	findOneAccountUsecase := accountUsecases.NewFindOneAccountUsecase(mockAccountRepo)
	createCardUsecase := cardUsecases.NewCreateCardUsecase(mockCardRepo, findOneAccountUsecase)
	findCardUsecase := cardUsecases.NewFindCardUsecase(mockCardRepo, findOneAccountUsecase)

	sut := NewCardHandler(createCardUsecase, findCardUsecase)

	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Status:   "active",
	}

	card := &domain.Card{
		Id:        1,
		Amount:    0,
		AccountId: 1,
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
		mockAccountRepo.On("FindById").Return(nil, errors.New("Internal server error"))
		defer mockAccountRepo.On("FindById").Unset()

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
		mockAccountRepo.On("FindById").Return(nil, nil)
		defer mockAccountRepo.On("FindById").Unset()

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

		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.Id),
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

		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("Save").Return(card, nil)
		defer mockCardRepo.On("Save").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.Id),
		}}

		sut.Create(c)

		var responseBody domain.Card
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, *card, responseBody)
	})

	t.Run("[FindCard] Error invalid card id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
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

		sut.FindCard(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid card id", responseBody["error"])
	})

	t.Run("[Create] Error account not found", func(t *testing.T) {
		mockAccountRepo.On("FindById").Return(nil, nil)
		defer mockAccountRepo.On("FindById").Unset()

		accountId := "1"

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
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

		sut.FindCard(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %s", accountId),
			responseBody["error"])
	})

	t.Run("[Create] Success to find card by id", func(t *testing.T) {
		account.Status = "active"

		mockAccountRepo.On("FindById").Return(account, nil)
		defer mockAccountRepo.On("FindById").Unset()

		mockCardRepo.On("FindById").Return(card, nil)
		defer mockCardRepo.On("FindById").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/card", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: fmt.Sprint(account.Id),
		},
			{
				Key:   "cardId",
				Value: fmt.Sprint(card.Id),
			}}

		sut.FindCard(c)

		var responseBody domain.Card
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, *card, responseBody)
	})

}
