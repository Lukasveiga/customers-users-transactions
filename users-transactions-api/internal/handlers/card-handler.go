package handlers

import (
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/tools"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	createCardUsecase *usecases.CreateCardUsecase
}

func NewCardHandler(createCardUsecase *usecases.CreateCardUsecase) *CardHandler {
	return &CardHandler{
		createCardUsecase: createCardUsecase,
	}
}

func (ch *CardHandler) Create(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	savedCard, err := ch.createCardUsecase.Create(tenantId, int32(accountId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}

		if ia, ok := err.(*shared.InactiveAccountError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": ia.Error()})
			return
		}

		tools.LogInternalServerError(c, "card handler", "Create", err)
		return
	}

	c.JSON(http.StatusCreated, savedCard)
}
