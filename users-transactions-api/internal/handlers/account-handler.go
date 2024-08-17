package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/dtos"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	createAccountUsecase *usecases.CreateAccountUsecase
}

func NewAccountHandler(createAccountUsecase *usecases.CreateAccountUsecase) *AccountHandler {
	return &AccountHandler{
		createAccountUsecase: createAccountUsecase,
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

func logInternalServerError(c *gin.Context, method string, err error) {
	slog.Error(
		"account handler",
		slog.String("method", method),
		slog.String("error", err.Error()),
	)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
