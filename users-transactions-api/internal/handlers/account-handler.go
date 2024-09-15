package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/handlers/dto"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/handlers/tools"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/usecases/account"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	createAccountUsecase   *usecases.CreateAccountUsecase
	findAllAccountsUsecase *usecases.FindAllAccountsUsecase
	findOneAccountUsecase  *usecases.FindOneAccountUsecase
	activeAccountUsecase   *usecases.ActiveAccountUsecase
	inactiveAccountUsecase *usecases.InactiveAccountUsecase
}

func NewAccountHandler(createAccountUsecase *usecases.CreateAccountUsecase, findAllAccountsUsecase *usecases.FindAllAccountsUsecase, findOneAccountUsecase *usecases.FindOneAccountUsecase, activeAccountUsecase *usecases.ActiveAccountUsecase, inactiveAccountUsecase *usecases.InactiveAccountUsecase) *AccountHandler {
	return &AccountHandler{
		createAccountUsecase:   createAccountUsecase,
		findAllAccountsUsecase: findAllAccountsUsecase,
		findOneAccountUsecase:  findOneAccountUsecase,
		activeAccountUsecase:   activeAccountUsecase,
		inactiveAccountUsecase: inactiveAccountUsecase,
	}
}

func (ah *AccountHandler) Create(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	savedAccount, err := ah.createAccountUsecase.Create(tenantId)

	if err != nil {
		tools.LogInternalServerError(c, "account handler", "Create", err)
		return
	}

	c.JSON(http.StatusCreated, dto.AccountToResponse(*savedAccount))
}

func (ah *AccountHandler) FindOne(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

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

		tools.LogInternalServerError(c, "account handler", "FindOne", err)
		return
	}

	c.JSON(http.StatusOK, dto.AccountToResponse(*account))
}

func (ah *AccountHandler) FindAll(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accounts, err := ah.findAllAccountsUsecase.FindAll(int32(tenantId))

	if err != nil {
		tools.LogInternalServerError(c, "account handler", "FindAll", err)
		return
	}

	accountsResponse := make([]dto.AccountResponse, 0)
	for _, account := range accounts {
		accountsResponse = append(accountsResponse, dto.AccountToResponse(account))
	}

	c.JSON(http.StatusOK, accountsResponse)
}

func (ah *AccountHandler) Active(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	err = ah.activeAccountUsecase.Active(tenantId, int32(accountId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": enf.Error()})
			return
		}

		tools.LogInternalServerError(c, "account handler", "Active", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account with id %d was actived successfully", accountId)})
}

func (ah *AccountHandler) Inactive(c *gin.Context) {
	tenantId, valid := tools.CheckTenantHeader(c)

	if !valid {
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	err = ah.inactiveAccountUsecase.Inactive(tenantId, int32(accountId))

	if err != nil {
		if enf, ok := err.(*shared.EntityNotFoundError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": enf.Error()})
			return
		}

		tools.LogInternalServerError(c, "account handler", "Inactive", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account with id %d was deleted successfully", accountId)})
}
