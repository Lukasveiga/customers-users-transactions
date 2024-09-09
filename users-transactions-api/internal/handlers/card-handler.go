package handlers

import (
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/dto"
	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/tools"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	createCardUsecase   *usecases.CreateCardUsecase
	findCardUsecase     *usecases.FindCardUsecase
	findAllCardsUsecase *usecases.FindAllCards
}

func NewCardHandler(createCardUsecase *usecases.CreateCardUsecase,
	findCardUsecase *usecases.FindCardUsecase, findAllCardsUsecase *usecases.FindAllCards) *CardHandler {
	return &CardHandler{
		createCardUsecase:   createCardUsecase,
		findCardUsecase:     findCardUsecase,
		findAllCardsUsecase: findAllCardsUsecase,
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

	c.JSON(http.StatusCreated, dto.CardToResponse(*savedCard))
}

func (ch *CardHandler) FindOne(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	cardId, err := strconv.ParseInt(c.Param("cardId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card id"})
		return
	}

	card, err := ch.findCardUsecase.FindOne(tenantId, int32(accountId), int32(cardId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}

		tools.LogInternalServerError(c, "card handler", "Create", err)
		return
	}

	c.JSON(http.StatusOK, dto.CardToResponse(*card))
}

func (ch *CardHandler) FindAll(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	cards, err := ch.findAllCardsUsecase.FindAll(tenantId, int32(accountId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}

		tools.LogInternalServerError(c, "card handler", "Create", err)
		return
	}

	cardsResponse := make([]dto.CardResponse, 0)
	for _, card := range cards {
		cardsResponse = append(cardsResponse, dto.CardToResponse(card))
	}

	c.JSON(http.StatusOK, cardsResponse)
}
