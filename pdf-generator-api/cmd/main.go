package main

import (
	"fmt"

	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/cmd/api"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/config"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/genproto"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/handlers"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/usecases"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	GRPC_PORT := config.GetEnv("GRPC_PORT")
	API_PORT := config.GetEnv("PORT")
	ENV := config.GetEnv("ENV")

	client := initClient(GRPC_PORT, ENV)
	initApiServer(client, API_PORT)
}

func initClient(PORT string, ENV string) genproto.TransactionInfoServiceClient {
	var host string

	switch ENV {
	case "prod":
		host = "nginx"
	default:
		host = "0.0.0.0"
	}

	serverAddress := fmt.Sprintf("%s:%s", host, PORT)
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
