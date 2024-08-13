package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
)

type TenantHandler struct {
	findOneTenantUsecase *usecases.FindOneTenantUseCase
}

func NewTenantHandler(findOneTenantUsecase *usecases.FindOneTenantUseCase) *TenantHandler {
	return &TenantHandler{
		findOneTenantUsecase: findOneTenantUsecase,
	}
}

func (th *TenantHandler) FindTenant(res http.ResponseWriter, req *http.Request, next http.Handler) {
	tenantId, err := strconv.ParseInt(req.Header.Get("tenant-id"), 0, 32)

	if err != nil {
		http.Error(res, "invalid tenant id", http.StatusBadRequest)
		return
	}

	_, err = th.findOneTenantUsecase.FindOne(int32(tenantId))

	if err != nil {
		if ae, ok := err.(*shared.EntityNotFoundError); ok {
			jsonData, err := json.Marshal(ae.Error())
			if err != nil {
				slog.Error(
					"tenant handler",
					slog.String("method", "FindTenant"),
					slog.String("error", err.Error()),
				)

				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.Error(res, string(jsonData), http.StatusBadRequest)
			return
		}

		slog.Error(
			"tenant handler",
			slog.String("method", "FindTenant"),
			slog.String("error", err.Error()),
		)

		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	next.ServeHTTP(res, req)
}
