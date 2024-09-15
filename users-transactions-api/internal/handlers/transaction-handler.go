package handlers

import (
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/handlers/dto"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/handlers/tools"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/transactions"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	createTransactionUsecase *usecases.CreateTransactionUsecase
	findTransactionUsecase   *usecases.FindTransactionUsecase
	findTransactionsUsecase  *usecases.FindAllTransactionsUsecase
}

func NewTransactionHandler(createTransactionUsecase *usecases.CreateTransactionUsecase,
	findTransactionUsecase *usecases.FindTransactionUsecase,
	findTransactionsUsecase *usecases.FindAllTransactionsUsecase) *TransactionHandler {
	return &TransactionHandler{
		createTransactionUsecase: createTransactionUsecase,
		findTransactionUsecase:   findTransactionUsecase,
		findTransactionsUsecase:  findTransactionsUsecase,
	}
}

func (th *TransactionHandler) Create(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	var request dto.TransactioRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := th.createTransactionUsecase.Create(tenantId, int32(accountId), dto.RequestToTransaction(request))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}

		if ve, ok := err.(*shared.ValidationError); ok {
			c.JSON(http.StatusNotFound, ve)
			return
		}

		tools.LogInternalServerError(c, "transaction handler", "Create", err)
		return
	}

	c.JSON(http.StatusCreated, dto.TransactionToResponse(*transaction))
}

func (th *TransactionHandler) FindOne(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	transactionId, err := strconv.ParseInt(c.Param("transactionId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	transaction, err := th.findTransactionUsecase.FindOne(tenantId, int32(accountId),
		int32(cardId), int32(transactionId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}
		tools.LogInternalServerError(c, "transaction handler", "Create", err)
		return
	}

	c.JSON(http.StatusOK, dto.TransactionToResponse(*transaction))
}

func (th *TransactionHandler) FindAll(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	transactions, err := th.findTransactionsUsecase.FindAll(tenantId, int32(accountId),
		int32(cardId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}
		tools.LogInternalServerError(c, "transaction handler", "Create", err)
		return
	}

	transactionsResponse := make([]dto.TransactionResponse, 0)

	for _, transaction := range transactions {
		transactionsResponse = append(transactionsResponse, dto.TransactionToResponse(transaction))
	}

	c.JSON(http.StatusOK, transactionsResponse)
}
