package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	usecases "github.com/Lukasveiga/customers-users-transactions/internal/usecases/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	client := genproto.NewTransactionInfoServiceClient(conn)

	result, err := usecases.SearchTransactionInformation(client, &genproto.Filter{
		TenantId:  1,
		AccountId: 1,
	})

	if err != nil {
		log.Fatal("error to search transaction: ", err)
	}

	for _, transInfo := range result {
		fmt.Printf("%v\n", transInfo)
	}
}
