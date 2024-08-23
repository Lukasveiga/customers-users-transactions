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
	findOneAccountUsecase  *usecases.FindOneAccountUsecase
	updateAccountUsecase   *usecases.UpdateAccountUsecase
	deleteAccountUsecase   *usecases.DeleteAccountUsecase
}

func NewAccountHandler(createAccountUsecase *usecases.CreateAccountUsecase, findAllAccountsUsecase *usecases.FindAllUsecase, findOneAccountUsecase *usecases.FindOneAccountUsecase, updateAccountUsecase *usecases.UpdateAccountUsecase, deleteAccountUsecase *usecases.DeleteAccountUsecase) *AccountHandler {
	return &AccountHandler{
		createAccountUsecase:   createAccountUsecase,
		findAllAccountsUsecase: findAllAccountsUsecase,
		findOneAccountUsecase:  findOneAccountUsecase,
		updateAccountUsecase:   updateAccountUsecase,
		deleteAccountUsecase:   deleteAccountUsecase,
	}
}

func (ah *AccountHandler) Create(c *gin.Context) {
	tenantId, valid := checkTenantHeader(c)

	if !valid {
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

		logInternalServerError(c, "Create", err)
		return
	}

	c.JSON(http.StatusCreated, savedAccount)
}

func (ah *AccountHandler) FindOne(c *gin.Context) {
	tenantId, valid := checkTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	account, err := ah.findOneAccountUsecase.FindOne(int32(tenantId), int32(accountId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": enf.Error()})
			return
		}

		logInternalServerError(c, "FindOne", err)
		return
	}

	c.JSON(http.StatusOK, account)
}

func (ah *AccountHandler) FindAll(c *gin.Context) {
	tenantId, valid := checkTenantHeader(c)

	if !valid {
		return
	}

	accounts, err := ah.findAllAccountsUsecase.FindAll(int32(tenantId))

	if err != nil {
		logInternalServerError(c, "FindAll", err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (ah *AccountHandler) Update(c *gin.Context) {
	tenantId, valid := checkTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	var accountDto *dtos.AccountDto

	if err := c.ShouldBindJSON(&accountDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Decoding Error"})
		return
	}

	account := accountDto.ToDomain()
	account.TenantId = int32(tenantId)

	updatedAccount, err := ah.updateAccountUsecase.Update(int32(accountId), account)

	if err != nil {
		if ae, ok := err.(*shared.AlreadyExistsError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": ae.Error()})
			return
		}

		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": enf.Error()})
			return
		}

		if ve, ok := err.(*shared.ValidationError); ok {
			c.JSON(http.StatusBadRequest, ve.Errors)
			return
		}

		logInternalServerError(c, "Update", err)
		return
	}

	c.JSON(http.StatusOK, updatedAccount)
}

func (ah *AccountHandler) Delete(c *gin.Context) {
	tenantId, valid := checkTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	err = ah.deleteAccountUsecase.Delete(tenantId, int32(accountId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": enf.Error()})
			return
		}

		logInternalServerError(c, "Update", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account with id %d was deleted successfully", accountId)})
}

func checkTenantHeader(c *gin.Context) (int32, bool) {
	tenantId, err := strconv.ParseInt(c.GetHeader("tenant-id"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant-id"})
		return -1, false
	}

	return int32(tenantId), true
}

func logInternalServerError(c *gin.Context, method string, err error) {
	slog.Error(
		"account handler",
		slog.String("method", method),
		slog.String("error", err.Error()),
	)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
