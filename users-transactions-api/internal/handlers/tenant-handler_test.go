package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTenantHandler(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockRepository)
	findOneTenantUsecase := usecases.NewFindOneTenantUseCase(mockRepo)
	sut := NewTenantHandler(findOneTenantUsecase)

	tenant := infra.Tenant{
		ID:   1,
		Name: "Tenant A",
	}

	t.Run("[FindTenant] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("tenant-id", "invalid-tenant-id")

		nextHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		}

		middleware := sut.FindTenant()
		middleware(c)
		nextHandler(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "invalid tenant id", responseBody["error"])
	})

	t.Run("[FindTenant] Tenant not found", func(t *testing.T) {
		mockRepo.On("GetTenant").Return(nil, sql.ErrNoRows)
		defer mockRepo.On("GetTenant").Unset()

		tenantId := "1"

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("tenant-id", tenantId)

		nextHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		}

		middleware := sut.FindTenant()
		middleware(c)
		nextHandler(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("tenant not found with id %s", tenantId),
			responseBody["error"])
	})

	t.Run("[FindTenant] Internal Server Error", func(t *testing.T) {
		internalErro := errors.New("internal server error")

		mockRepo.On("GetTenant").Return(nil, internalErro)
		defer mockRepo.On("GetTenant").Unset()

		tenantId := "1"

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("tenant-id", tenantId)

		nextHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		}

		middleware := sut.FindTenant()
		middleware(c)
		nextHandler(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error", responseBody["error"])
	})

	t.Run("[FindTenant] Success", func(t *testing.T) {
		mockRepo.On("GetTenant").Return(tenant, nil)
		defer mockRepo.On("GetTenant").Unset()

		tenantId := strconv.Itoa(int(tenant.ID))

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("tenant-id", tenantId)

		nextHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		}

		middleware := sut.FindTenant()
		middleware(c)
		nextHandler(c)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "\"\"", res.Body.String())
	})
}
