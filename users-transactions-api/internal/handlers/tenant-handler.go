package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	findOneTenantUsecase *usecases.FindOneTenantUseCase
}

func NewTenantHandler(findOneTenantUsecase *usecases.FindOneTenantUseCase) *TenantHandler {
	return &TenantHandler{
		findOneTenantUsecase: findOneTenantUsecase,
	}
}

func (th *TenantHandler) FindTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, err := strconv.ParseInt(c.GetHeader("tenant-id"), 0, 32)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
			c.Abort()
			return
		}

		_, err = th.findOneTenantUsecase.FindOne(int32(tenantId))

		if err != nil {
			if ae, ok := err.(*shared.EntityNotFoundError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": ae.Error()})
				c.Abort()
				return
			}

			slog.Error(
				"tenant handler",
				slog.String("method", "FindTenant"),
				slog.String("error", err.Error()),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		c.Next()
	}
}
