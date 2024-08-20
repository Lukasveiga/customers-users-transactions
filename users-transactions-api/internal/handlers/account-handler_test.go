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
	findAllAccountsUsecase := usecases.NewFindAllAccountsUsecase(mockRepo)
	sut := NewAccountHandler(accountCreateUsecase, findAllAccountsUsecase)

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

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "{\"error\":\"Decoding Error\"}", res.Body.String()) // TODO: improve the comparison
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

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("{\"error\":\"account already exists with id %s\"}", account.Number),
			res.Body.String())
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

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "{\"error\":\"Internal Server Error\"}", res.Body.String())
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
}
