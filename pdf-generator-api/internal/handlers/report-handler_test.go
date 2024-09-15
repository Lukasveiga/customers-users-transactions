package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/mocks"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReportHandler(t *testing.T) {
	t.Parallel()

	mockPdfGenerator := new(mocks.MockPdfGenerator)
	transReport := usecases.NewTransactionReport(client, mockPdfGenerator)
	sut := NewReportHandler(transReport)

	t.Run("[SendReport] Invalid tenant id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transations-report", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{
			Key:   "tenantId",
			Value: "invalid",
		}}

		sut.SendReport(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid tenant id", responseBody["error"])
	})

	t.Run("[SendReport] Invalid account id", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transations-report", nil)
		c.Request.Header.Set("Content-Type", "application/pdf")
		c.Params = []gin.Param{
			{
				Key:   "tenantId",
				Value: "1",
			},
			{
				Key:   "accountId",
				Value: "invalid",
			}}

		sut.SendReport(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "Invalid account id", responseBody["error"])
	})

	t.Run("[SendReport] Transaction not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/transations-report", nil)
		c.Request.Header.Set("Content-Type", "application/pdf")
		c.Params = []gin.Param{
			{
				Key:   "tenantId",
				Value: "1",
			},
			{
				Key:   "accountId",
				Value: "99",
			}}

		sut.SendReport(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})

	t.Run("[SendReport] File not found", func(t *testing.T) {
		path := "./invalid-path/test.pdf"

		mockPdfGenerator.On("Generate").Return(path, nil)
		defer mockPdfGenerator.On("Generate").Unset()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)

		c.Request = httptest.NewRequest("GET", "/account", nil)
		c.Request.Header.Set("Content-Type", "application/pdf")
		c.Params = []gin.Param{
			{
				Key:   "tenantId",
				Value: "1",
			},
			{
				Key:   "accountId",
				Value: "1",
			}}

		sut.SendReport(c)

		var responseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&responseBody)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, "File not found", responseBody["error"])
	})
}
