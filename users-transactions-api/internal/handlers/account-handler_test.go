package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/dtos"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccountHandler(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockAccountRepository)
	accountCreateUsecase := usecases.NewCreateAccountUsecase(mockRepo)
	findOneAccountUsecase := usecases.NewFindOneAccountUsecase(mockRepo)
	findAllAccountsUsecase := usecases.NewFindAllAccountsUsecase(mockRepo)
	updateAccountUsecase := usecases.NewUpdateAccountUsecase(mockRepo)
	deleteAccountUSecase := usecases.NewDeleteAccountUsecase(mockRepo)
	sut := NewAccountHandler(accountCreateUsecase, findAllAccountsUsecase, findOneAccountUsecase, updateAccountUsecase, deleteAccountUSecase)

	number := uuid.New().String()

	account := &domain.Account{
		Id:       1,
		TenantId: 1,
		Number:   number,
		Status:   "active",
	}

	accountDto := &dtos.AccountDto{
		Number: number,
		Status: "active",
	}

	t.Run("[Create] Decoding error bad request", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Decoding Error", responseBody["error"])
	})

	t.Run("[Create] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[Create] Account already exists error bad request", func(t *testing.T) {
		mockRepo.On("FindByNumber").Return(account, nil)
		defer mockRepo.On("FindByNumber").Unset()

		body, err := json.Marshal(accountDto)

		assert.NoError(t, err)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody map[string]string
		err = json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account already exists with id %s", account.Number),
			responseBody["error"])
	})

	t.Run("[Create] Invalid input error bad request", func(t *testing.T) {
		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		invalidAccountDto := &dtos.AccountDto{
			Number: "1",
			Status: "invalid",
		}

		body, err := json.Marshal(invalidAccountDto)

		assert.NoError(t, err)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody map[string]string
		err = json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "must be a valid uuid", responseBody["number"])
		assert.Equal(t, "must be active or inactive", responseBody["status"])
	})

	t.Run("[Create] Internal server error", func(t *testing.T) {
		mockRepo.On("FindByNumber").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("FindByNumber").Unset()

		body, err := json.Marshal(accountDto)

		assert.NoError(t, err)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody map[string]string
		err = json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[Create] Account created successfully", func(t *testing.T) {
		mockRepo.On("FindByNumber").Return(nil, nil)
		defer mockRepo.On("FindByNumber").Unset()

		mockRepo.On("Create").Return(account, nil)
		defer mockRepo.On("Create").Unset()

		body, err := json.Marshal(accountDto)

		assert.NoError(t, err)

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody domain.Account
		err = json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, *account, responseBody)
	})

	t.Run("[FindOne] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[FindOne] Invalid account id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "invalid",
		}}

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid account id", responseBody["error"])
	})

	t.Run("[FindOne] Account not found by id", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, nil)
		defer mockRepo.On("FindById").Unset()

		accountId := "1"

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: accountId,
		}}

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %s", accountId), responseBody["error"])
	})

	t.Run("[FindOne] Internal server error", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("FindById").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.FindOne(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[FindOne] Success", func(t *testing.T) {
		mockRepo.On("FindById").Return(account, nil)
		defer mockRepo.On("FindById").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.FindOne(c)

		var responseBody domain.Account
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, *account, responseBody)
	})

	t.Run("[FindAll] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[FindAll] Internal server error", func(t *testing.T) {
		mockRepo.On("FindAll").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("FindAll").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.FindAll(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[FindAll] Success with empty array", func(t *testing.T) {
		accounts := make([]*domain.Account, 0)

		mockRepo.On("FindAll").Return(accounts, nil)
		defer mockRepo.On("FindAll").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.FindAll(c)

		var responseBody []*domain.Account
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, accounts, responseBody)
	})

	t.Run("[FindAll] Success", func(t *testing.T) {
		accounts := make([]*domain.Account, 0)
		accounts = append(accounts, account)

		mockRepo.On("FindAll").Return(accounts, nil)
		defer mockRepo.On("FindAll").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.FindAll(c)

		var responseBody []*domain.Account
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Len(t, responseBody, 1)
		assert.Equal(t, accounts, responseBody)
	})
}
