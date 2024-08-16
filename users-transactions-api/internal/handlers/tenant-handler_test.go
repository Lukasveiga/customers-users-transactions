package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	"github.com/Lukasveiga/customers-users-transaction/internal/mocks"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
	"github.com/stretchr/testify/assert"
)

func TestTenantHandler(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.MockTenantRepository)
	findOneTenantUsecase := usecases.NewFindOneTenantUseCase(mockRepo)
	sut := NewTenantHandler(findOneTenantUsecase)

	tenant := &domain.Tenant{
		Id:   int32(1),
		Name: "Tenant A",
	}

	t.Run("[FindTenant] Invalid tenant id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("tenant-id", "invalid-tenant-id")

		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		sut.FindTenant(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "invalid tenant id\n", res.Body.String())
	})

	t.Run("[FindTenant] Tenant not found", func(t *testing.T) {
		mockRepo.On("FindById").Return(nil, nil)
		defer mockRepo.On("FindById").Unset()

		tenantId := "1"

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("tenant-id", tenantId)

		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		sut.FindTenant(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("\"tenant not found with id %s\"\n", tenantId),
			res.Body.String())
	})

	t.Run("[FindTenant] Internal Server Error", func(t *testing.T) {
		internalErro := errors.New("internal server error")

		mockRepo.On("FindById").Return(nil, internalErro)
		defer mockRepo.On("FindById").Unset()

		tenantId := "1"

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("tenant-id", tenantId)

		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		sut.FindTenant(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "Internal Server Error\n", res.Body.String())
	})

	t.Run("[FindTenant] Success", func(t *testing.T) {
		mockRepo.On("FindById").Return(tenant, nil)
		defer mockRepo.On("FindById").Unset()

		tenantId := strconv.Itoa(int(tenant.Id))

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("tenant-id", tenantId)

		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		sut.FindTenant(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "", res.Body.String())
	})
}
