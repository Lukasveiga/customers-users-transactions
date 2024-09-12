package api

import (
	"github.com/Lukasveiga/customers-users-transactions/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Routes(handler *handlers.ReportHandler) *gin.Engine {
	baseUrl := "/api/v1"
	router := gin.Default()

	router.GET(baseUrl+"/accounts/:accountId/tenant/:tenantId/transactions.pdf", handler.SendReport)

	return router
}
