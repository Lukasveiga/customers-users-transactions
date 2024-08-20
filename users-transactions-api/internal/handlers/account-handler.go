package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/dtos"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	createAccountUsecase   *usecases.CreateAccountUsecase
	findAllAccountsUsecase *usecases.FindAllUsecase
}

func NewAccountHandler(createAccountUsecase *usecases.CreateAccountUsecase, findAllAccountsUsecase *usecases.FindAllUsecase) *AccountHandler {
	return &AccountHandler{
		createAccountUsecase:   createAccountUsecase,
		findAllAccountsUsecase: findAllAccountsUsecase,
	}
}

func (ah *AccountHandler) Create(c *gin.Context) {
	tenantId, err := strconv.ParseInt(c.GetHeader("tenant-id"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant-id"})
		return
	}

	var accountDto *dtos.AccountDto

	if err := c.ShouldBindJSON(&accountDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Decoding Error"})
		return
	}

	account := accountDto.ToDomain()
	account.TenantId = int32(tenantId)

	savedAccount, err := ah.createAccountUsecase.Create(account)

	if err != nil {
		if ae, ok := err.(*shared.AlreadyExistsError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": ae.Error()})
			return
		}

		if ve, ok := err.(*shared.ValidationError); ok {
			c.JSON(http.StatusBadRequest, ve.Errors)
			return
		}

		logInternalServerError(c, "create", err)
		return
	}

	c.JSON(http.StatusCreated, savedAccount)
}

func (ah *AccountHandler) FindAll(c *gin.Context) {
	tenantId, err := strconv.ParseInt(c.GetHeader("tenant-id"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant-id"})
		return
	}

	accounts, err := ah.findAllAccountsUsecase.FindAll(int32(tenantId))

	if err != nil {
		logInternalServerError(c, "findAll", err)
		return
	}

	fmt.Println(accounts)

	c.JSON(http.StatusOK, accounts)
}

func logInternalServerError(c *gin.Context, method string, err error) {
	slog.Error(
		"account handler",
		slog.String("method", method),
		slog.String("error", err.Error()),
	)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
