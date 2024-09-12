package usecases

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/genproto"
	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	_, err := ti.repo.GetTenant(stream.Context(), int32(filter.GetTenantId()))

	if err != nil {
		if err == sql.ErrNoRows {
			return status.Errorf(codes.NotFound, fmt.Sprintf("tenant with id %d not found", filter.GetTenantId()), err)
		}

		slog.Error(
			"error to find tenant by id",
			slog.String("err", err.Error()),
		)
		return status.Errorf(codes.Internal, "Internal repository error: %v", err)
	}

	_, err = ti.repo.GetAccount(stream.Context(), infra.GetAccountParams{
		TenantID: int32(filter.GetTenantId()),
		ID:       int32(filter.GetAccountId()),
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return status.Errorf(codes.NotFound, fmt.Sprintf("account with id %d not found", filter.GetAccountId()), err)
		}

		slog.Error(
			"error to find account by id",
			slog.String("err", err.Error()),
		)
		return status.Errorf(codes.Internal, "Internal repository error: %v", err)
	}

	result, err := ti.repo.SearchTransactions(stream.Context(), infra.SearchTransactionsParams{
		TenantID:  int32(filter.TenantId),
		Accountid: int32(filter.AccountId),
	})

	if err != nil {
		slog.Error(
			"error search transactions information",
			slog.String("err", err.Error()),
		)
		return status.Errorf(codes.Internal, "Internal repository error: %v", err)
	}

	for _, t := range result {
		response := &genproto.TransactionInfo{
			AccountId: filter.GetAccountId(),
			Kind:      t.Kind,
			Value:     float64(t.Value),
		}

		err := stream.Send(&genproto.SearchTransactionInfoResponse{TransactionInfo: response})

		if err != nil {
			slog.Error(
				"error to send response stream",
				slog.String("err", err.Error()),
			)
			return status.Errorf(codes.Internal, "Cannot send response stream: %v", err)
		}
	}
	return nil
}
