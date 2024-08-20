package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/cmd/api/factory"
	"github.com/Lukasveiga/customers-users-transaction/cmd/api/router"
	"github.com/Lukasveiga/customers-users-transaction/config"
)

func main() {
	PORT := config.GetEnv("PORT")
	ENV := config.GetEnv("ENV")

	logConfig := &config.LoggerConfig{
		Env:     ENV,
		LogPath: "./tmp/logs.log",
	}

	config.InitLogger(logConfig)

	psqlInfo := factory.GetDbUrlConn(ENV)

	dbConnection := initDbConnection(psqlInfo)

	startServer(PORT, dbConnection)
}

func startServer(PORT string, dbConnection *sql.DB) {
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

func initDbConnection(psqlInfo string) *sql.DB {
	slog.Info("database connection established")
	return config.InitConfig(psqlInfo)
}
