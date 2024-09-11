package usecases

import (
	"github.com/Lukasveiga/customers-users-transaction/internal/genproto"
	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
)

type TransactionInfo struct {
	repo infra.Querier
}

func NewTransactionInfo(repo infra.Querier) *TransactionInfo {
	return &TransactionInfo{
		repo: repo,
	}
}

func (ti *TransactionInfo) SearchTransactionInfo(req *genproto.SearchTransactionInfoRequest,
	stream genproto.TransactionInfoService_SearchTransactionInfoServer) error {

	filter := req.GetFilter()

	result, err := ti.repo.SearchTransactions(stream.Context(), infra.SearchTransactionsParams{
		TenantID:  int32(filter.TenantId),
		Accountid: int32(filter.AccountId),
	})

	if err != nil {
		return err
	}

	for _, t := range result {
		response := &genproto.TransactionInfo{
			AccountId: filter.GetAccountId(),
			Kind:      t.Kind,
			Value:     float64(t.Value),
		}

		err := stream.Send(&genproto.SearchTransactionInfoResponse{TransactionInfo: response})

		if err != nil {
			return err
		}
	}
	return nil
}
