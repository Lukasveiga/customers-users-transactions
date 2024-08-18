package router

import (
	"net/http"

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

	router.GET(baseUrl+"/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	router.Use(handlers.TenantHandler.FindTenant())

	account := router.Group(baseUrl)
	{
		account.POST("/account", handlers.AccountHandler.Create)
	}

	return router
}
