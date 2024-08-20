package router

import (
	"net/http"

	"github.com/Lukasveiga/customers-users-transaction/cmd/api/factory"
	"github.com/gin-gonic/gin"
)

func Routes(handlers *factory.Handlers) *gin.Engine {
	baseUrl := "/api/v1"
	router := gin.Default()

	router.GET(baseUrl+"/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	router.Use(handlers.TenantHandler.FindTenant())

	account := router.Group(baseUrl)
	{
		account.POST("/account", handlers.AccountHandler.Create)
		account.GET("/account", handlers.AccountHandler.FindAll)
	}

	return router
}
