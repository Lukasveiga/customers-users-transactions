package usecases

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
)

func SearchTransactionInformation(client genproto.TransactionInfoServiceClient,
	filter *genproto.Filter) ([]*genproto.TransactionInfo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &genproto.SearchTransactionInfoRequest{Filter: filter}

	stream, err := client.SearchTransactionInfo(ctx, req)

	if err != nil {
		log.Fatal("cannot search transaction information: ", err)
	}

	transactionInfo := make([]*genproto.TransactionInfo, 0)

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			return transactionInfo, nil
		}

		if err != nil {
			log.Fatal("cannot receive transaction information response: ", err)
		}

		transactionInfo = append(transactionInfo, res.GetTransactionInfo())
	}
}
