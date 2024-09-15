package router

import (
	"net/http"

	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/cmd/api/factory"
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
		account.GET("/account/:accountId", handlers.AccountHandler.FindOne)
		account.GET("/account", handlers.AccountHandler.FindAll)
		account.PUT("/account/:accountId", handlers.AccountHandler.Active)
		account.DELETE("/account/:accountId", handlers.AccountHandler.Inactive)
	}

	card := router.Group(baseUrl)
	{
		card.POST("/card/:accountId", handlers.CardHandler.Create)
		card.GET("/card/:cardId/account/:accountId", handlers.CardHandler.FindOne)
		card.GET("/card/account/:accountId", handlers.CardHandler.FindAll)
	}

	transaction := router.Group(baseUrl)
	{
		transaction.POST("/transaction/account/:accountId", handlers.TransactionHandler.Create)
		transaction.GET("/transaction/:transactionId/account/:accountId/card/:cardId",
			handlers.TransactionHandler.FindOne)
		transaction.GET("/transaction/account/:accountId/card/:cardId",
			handlers.TransactionHandler.FindAll)
	}
	return router
}
