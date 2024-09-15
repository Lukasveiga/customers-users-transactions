package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/handlers/dto"
	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/account"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAccountHandler(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	accountCreateUsecase := usecases.NewCreateAccountUsecase(mockRepo)
	findOneAccountUsecase := usecases.NewFindOneAccountUsecase(mockRepo)
	findAllAccountsUsecase := usecases.NewFindAllAccountsUsecase(mockRepo)
	updateAccountUsecase := usecases.NewActiveAccountUsecase(mockRepo)
	deleteAccountUSecase := usecases.NewInactiveAccountUsecase(mockRepo)

	sut := NewAccountHandler(accountCreateUsecase, findAllAccountsUsecase, findOneAccountUsecase, updateAccountUsecase, deleteAccountUSecase)

	account := infra.Account{
		ID:       1,
		TenantID: 1,
		Status:   "active",
	}

	t.Run("[Create] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[Create] Internal server error", func(t *testing.T) {
		mockRepo.On("CreateAccount").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("CreateAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[Create] Account created successfully", func(t *testing.T) {
		mockRepo.On("CreateAccount").Return(account, nil)
		defer mockRepo.On("CreateAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("POST", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.Create(c)

		var responseBody dto.AccountResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, dto.AccountToResponse(account), responseBody)
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
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

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
		mockRepo.On("GetAccount").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("GetAccount").Unset()

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
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

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

		var responseBody dto.AccountResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, dto.AccountToResponse(account), responseBody)
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
		mockRepo.On("GetAccounts").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("GetAccounts").Unset()

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
		accounts := make([]infra.Account, 0)

		mockRepo.On("GetAccounts").Return(accounts, nil)
		defer mockRepo.On("GetAccounts").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.FindAll(c)

		var responseBody []infra.Account
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, accounts, responseBody)
	})

	t.Run("[FindAll] Success", func(t *testing.T) {
		accounts := make([]infra.Account, 0)
		accounts = append(accounts, account)

		mockRepo.On("GetAccounts").Return(accounts, nil)
		defer mockRepo.On("GetAccounts").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")

		sut.FindAll(c)

		var responseBody []dto.AccountResponse
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Len(t, responseBody, 1)
		assert.Equal(t, []dto.AccountResponse{dto.AccountToResponse(account)}, responseBody)
	})

	t.Run("[Active] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("PUT", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.Active(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[Active] Invalid account Id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("PUT", "/account", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "invalid",
		}}

		sut.Active(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid account id", responseBody["error"])
	})

	t.Run("[Active] Account not found by id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		accountId := "1"

		c.Request = httptest.NewRequest("PUT", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: accountId,
		}}

		sut.Active(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %s", accountId),
			responseBody["error"])
	})

	t.Run("[Active] Internal server error", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, errors.New("Internal server error"))
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("PUT", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.Active(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[Active] Account updated successfully", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("UpdateAccount").Return(account, nil)
		defer mockRepo.On("UpdateAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("PUT", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.Active(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("Account with id %d was actived successfully", account.ID),
			responseBody["message"])
	})

	t.Run("[Inactive] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("DELETE", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "invalid")

		sut.Inactive(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant-id", responseBody["error"])
	})

	t.Run("[Inactive] Invalid account id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("DELETE", "/account", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "invalid",
		}}

		sut.Inactive(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid account id", responseBody["error"])
	})

	t.Run("[Inactive] Account not found by id", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		accountId := "1"

		c.Request = httptest.NewRequest("DELETE", "/account", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: accountId,
		}}

		sut.Inactive(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("account not found with id %s", accountId),
			responseBody["error"])
	})

	t.Run("[Inactive] Internal server error", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(nil, errors.New("Internal server error"))
		defer mockRepo.On("GetAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("DELETE", "/account", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.Inactive(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[Inactive] Account deleted successfully", func(t *testing.T) {
		mockRepo.On("GetAccount").Return(account, nil)
		defer mockRepo.On("GetAccount").Unset()

		mockRepo.On("UpdateAccount").Return(account, nil)
		defer mockRepo.On("UpdateAccount").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("DELETE", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("tenant-id", "1")
		c.Params = []gin.Param{{
			Key:   "accountId",
			Value: "1",
		}}

		sut.Inactive(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("Account with id %d was deleted successfully", account.ID),
			responseBody["message"])
	})
}
