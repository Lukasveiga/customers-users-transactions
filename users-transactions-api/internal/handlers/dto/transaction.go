package dto

import infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"

type TransactioRequest struct {
	CardId int32  `json:"card_id"`
	Kind   string `json:"kind"`
	Value  int64  `json:"value"`
}

type TransactionResponse struct {
	ID     int32  `json:"id"`
	CardId int32  `json:"card_id"`
	Kind   string `json:"kind"`
	Value  int64  `json:"value"`
}

func TransactionToResponse(transaction infra.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:     transaction.ID,
		CardId: transaction.CardID,
		Kind:   transaction.Kind,
		Value:  transaction.Value,
	}
}

func RequestToTransaction(request TransactioRequest) infra.Transaction {
	return infra.Transaction{
		CardID: request.CardId,
		Kind:   request.Kind,
		Value:  request.Value,
	}
}
