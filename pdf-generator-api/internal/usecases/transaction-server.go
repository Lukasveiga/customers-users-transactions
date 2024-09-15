package usecases

import (
	"database/sql"
	"log"

	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionInfoServer struct {
	transactionInfoRepo TransactionInfoRepository
}

func NewTransactionInfoServer(transactionInfoRepo TransactionInfoRepository) *TransactionInfoServer {
	return &TransactionInfoServer{transactionInfoRepo: transactionInfoRepo}
}

func (server *TransactionInfoServer) SearchTransactionInfo(req *genproto.SearchTransactionInfoRequest,
	stream genproto.TransactionInfoService_SearchTransactionInfoServer) error {

	filter := req.GetFilter()

	err := server.transactionInfoRepo.Search(stream.Context(), filter,
		func(transInfo *genproto.TransactionInfo) error {
			res := &genproto.SearchTransactionInfoResponse{TransactionInfo: transInfo}

			err := stream.Send(res)

			if err != nil {
				return err
			}

			log.Printf("sent transaction info for account with id: %d", transInfo.GetAccountId())
			return nil
		})

	if err != nil {
		if err == sql.ErrNoRows {
			return status.Errorf(codes.NotFound, "sql not found err: %v", err)
		}
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}
