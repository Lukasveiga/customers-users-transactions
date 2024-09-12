package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"

	"github.com/Lukasveiga/customers-users-transaction/cmd/api/factory"
	"github.com/Lukasveiga/customers-users-transaction/cmd/api/router"
	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/genproto"
	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	usecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/grpc"
	"google.golang.org/grpc"
)

func main() {
	PORT := config.GetEnv("PORT")
	ENV := config.GetEnv("ENV")
	GRPC_PORT := config.GetEnv("GRPC_PORT")

	logConfig := &config.LoggerConfig{
		Env:     ENV,
		LogPath: "./tmp/logs.log",
	}

	config.InitLogger(logConfig)

	psqlInfo := factory.GetDbUrlConn(ENV)

	dbConnection := initDbConnection(psqlInfo)

	go startGrpcServer(GRPC_PORT, dbConnection)

	startApiServer(PORT, dbConnection)
}

func startApiServer(PORT string, dbConnection *sql.DB) {
	handlers := factory.InitHandlers(dbConnection)

	router := router.Routes(handlers)

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))

	if err != nil {
		slog.Error("cannot start server",
			slog.String("error", err.Error()),
		)

		panic(err)
	}
}

func startGrpcServer(PORT string, dbConnection *sql.DB) {
	repo := infra.New(dbConnection)
	server := usecases.NewTransactionInfo(repo)

	grpcServer := grpc.NewServer()
	genproto.RegisterTransactionInfoServiceServer(grpcServer, server)

	address := fmt.Sprintf("0.0.0.0:%s", PORT)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		slog.Error("cannot start server",
			slog.String("error", err.Error()),
		)
		panic(err)
	}

	err = grpcServer.Serve(listener)

	if err != nil {
		slog.Error("cannot start server",
			slog.String("error", err.Error()),
		)
		panic(err)
	}
}

func initDbConnection(psqlInfo string) *sql.DB {
	slog.Info("database connection established")
	return config.InitConfig(psqlInfo)
}
