package dto

import infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"

type CardResponse struct {
	ID        int32 `json:"id"`
	AccountID int32 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func CardToResponse(card infra.Card) CardResponse {
	return CardResponse{
		ID:        card.ID,
		AccountID: card.AccountID,
		Amount:    card.Amount,
	}
}
