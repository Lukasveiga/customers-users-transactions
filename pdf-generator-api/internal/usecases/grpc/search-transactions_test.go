package usecases

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientSeachTransactionInfo(t *testing.T) {
	t.Parallel()

	filter := &genproto.Filter{
		TenantId:  1,
		AccountId: 1,
	}

	transInfoRepo := NewInMemoryTransactionInfoRepository()

	serverAddress := startServer(t, transInfoRepo)
	transInfoClient := newTestTransactionInfoClient(t, serverAddress)

	req := &genproto.SearchTransactionInfoRequest{Filter: filter}
	stream, err := transInfoClient.SearchTransactionInfo(context.Background(), req)
	assert.NoError(t, err)

	found := 0

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		assert.NoError(t, err)
		assert.Equal(t, filter.GetAccountId(), res.GetTransactionInfo().GetAccountId())

		found += 1
	}

	assert.Equal(t, 2, found)
}

func startServer(t *testing.T, transRepo TransactionInfoRepository) string {
	transInfoServer := NewTransactionInfoServer(transRepo)

	grpcServer := grpc.NewServer()
	genproto.RegisterTransactionInfoServiceServer(grpcServer, transInfoServer)

	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)

	go func() {
		err := grpcServer.Serve(listener)
		assert.NoError(t, err)
	}()

	return listener.Addr().String()
}

func newTestTransactionInfoClient(t *testing.T, serverAddress string) genproto.TransactionInfoServiceClient {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	return genproto.NewTransactionInfoServiceClient(conn)
}
