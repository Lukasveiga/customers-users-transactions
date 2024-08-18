package router

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AccountHandler *handlers.AccountHandler
	TenantHandler  *handlers.TenantHandler
}

func Routes(handlers *Handlers) *gin.Engine {
	baseUrl := "/api/v1"
	router := gin.Default()

	router.Use(handlers.TenantHandler.FindTenant())

	account := router.Group(baseUrl)
	{
		account.POST("/account", handlers.AccountHandler.Create)
	}

	return router
}
