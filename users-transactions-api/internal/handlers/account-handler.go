package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Lukasveiga/customers-users-transaction/internal/handlers/dtos"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
)

type AccountHandler struct {
	createAccountUsecase *usecases.CreateAccountUsecase
}

func NewAccountHandler(createAccountUsecase *usecases.CreateAccountUsecase) *AccountHandler {
	return &AccountHandler{
		createAccountUsecase: createAccountUsecase,
	}
}

func (ah AccountHandler) Create(res http.ResponseWriter, req *http.Request) {
	tenantId, _ := strconv.ParseInt(req.Header.Get("tenant-id"), 0, 32)
	var accountDto *dtos.AccountDto

	err := json.NewDecoder(req.Body).Decode(&accountDto)

	if err != nil {
		http.Error(res, "Decoding Error", http.StatusBadRequest)
		return
	}

	account := accountDto.ToDomain()
	account.TenantId = int32(tenantId)

	savedAccount, err := ah.createAccountUsecase.Create(account)

	if err != nil {
		if ae, ok := err.(*shared.AlreadyExistsError); ok {
			jsonData, err := json.Marshal(ae.Error())
			if err != nil {
				logInternalServerError(res, "create", err)
				return
			}

			http.Error(res, string(jsonData), http.StatusBadRequest)
			return
		}

		if ve, ok := err.(*shared.ValidationError); ok {
			jsonData, err := json.Marshal(ve.Errors)
			if err != nil {
				logInternalServerError(res, "create", err)
				return
			}

			http.Error(res, string(jsonData), http.StatusBadRequest)
			return
		}

		logInternalServerError(res, "create", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(res).Encode(savedAccount)

	if err != nil {
		logInternalServerError(res, "create", err)
		return
	}
}

func logInternalServerError(res http.ResponseWriter, method string, err error) {
	slog.Error(
		"account handler",
		slog.String("method", method),
		slog.String("error", err.Error()),
	)

	http.Error(res, "Internal Server Error", http.StatusInternalServerError)
}
