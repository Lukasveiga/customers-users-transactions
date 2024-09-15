package usecases

import (
	"context"
	"io"
	"time"

	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/genproto"
)

func SearchTransactionInformation(client genproto.TransactionInfoServiceClient,
	filter *genproto.Filter) ([]*genproto.TransactionInfo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &genproto.SearchTransactionInfoRequest{Filter: filter}

	stream, err := client.SearchTransactionInfo(ctx, req)

	if err != nil {
		return nil, err
	}

	transactionInfo := make([]*genproto.TransactionInfo, 0)

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			return transactionInfo, nil
		}

		if err != nil {
			return nil, err
		}

		transactionInfo = append(transactionInfo, res.GetTransactionInfo())
	}
}
