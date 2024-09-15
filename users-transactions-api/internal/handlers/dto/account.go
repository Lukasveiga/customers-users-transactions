package dto

import infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"

type AccountResponse struct {
	ID       int32  `json:"id"`
	TenantID int32  `json:"tenant_id"`
	Status   string `json:"status"`
}

func AccountToResponse(account infra.Account) AccountResponse {
	return AccountResponse{
		ID:       account.ID,
		TenantID: account.TenantID,
		Status:   account.Status,
	}
}
