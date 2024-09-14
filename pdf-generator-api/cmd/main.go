package main

import (
	"fmt"

	"github.com/Lukasveiga/customers-users-transactions/cmd/api"
	"github.com/Lukasveiga/customers-users-transactions/config"
	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	"github.com/Lukasveiga/customers-users-transactions/internal/handlers"
	"github.com/Lukasveiga/customers-users-transactions/internal/usecases"
	"github.com/Lukasveiga/customers-users-transactions/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	GRPC_PORT := config.GetEnv("GRPC_PORT")
	API_PORT := config.GetEnv("API_PORT")

	client := initClient(GRPC_PORT)
	initApiServer(client, API_PORT)
}

func initClient(PORT string) genproto.TransactionInfoServiceClient {
	serverAddress := fmt.Sprintf("0.0.0.0:%s", PORT)
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return genproto.NewTransactionInfoServiceClient(conn)
}

func initApiServer(client genproto.TransactionInfoServiceClient, PORT string) {
	transactionReport := usecases.NewTransactionReport(client, utils.NewGofpdfGenerator())
	reportHandler := handlers.NewReportHandler(transactionReport)

	router := api.Routes(reportHandler)

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))

	if err != nil {
		panic(err)
	}
}
