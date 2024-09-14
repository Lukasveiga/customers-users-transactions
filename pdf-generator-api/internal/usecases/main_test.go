package usecases

import (
	"net"
	"os"
	"testing"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client genproto.TransactionInfoServiceClient
)

func TestMain(m *testing.M) {
	transInfoRepo := NewInMemoryTransactionInfoRepository()

	transInfoServer := NewTransactionInfoServer(transInfoRepo)

	grpcServer := grpc.NewServer()
	genproto.RegisterTransactionInfoServiceServer(grpcServer, transInfoServer)

	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		panic(err)
	}

	go func() {
		err := grpcServer.Serve(listener)

		if err != nil {
			panic(err)
		}
	}()

	serverAddress := listener.Addr().String()

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	client = genproto.NewTransactionInfoServiceClient(conn)

	os.Exit(m.Run())
}
