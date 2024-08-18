package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/cmd/api/router"
	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/handlers"
	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	tenantUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
)

func main() {
	PORT := config.GetEnv("PORT")
	ENV := config.GetEnv("ENV")

	config.InitLogger(&config.LoggerConfig{
		Env:     ENV,
		LogPath: "./logs.log",
	})

	var (
		host, db_port, user, password, dbname string
	)

	switch ENV {
	case "prod":
		host = config.GetEnv("DB_HOST")
		db_port = config.GetEnv("DB_PORT")
		user = config.GetEnv("DB_USERNAME")
		password = config.GetEnv("DB_PASSWORD")
		dbname = config.GetEnv("DB_NAME")
	default:
		host = config.GetEnv("DB_HOST_DEV")
		db_port = config.GetEnv("DB_PORT_DEV")
		user = config.GetEnv("DB_USERNAME_DEV")
		password = config.GetEnv("DB_PASSWORD_DEV")
		dbname = config.GetEnv("DB_NAME_DEV")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, db_port, user, password, dbname)

	dbConnection := initDbConnection(psqlInfo)

	accountPgRepository := infra.NewPgAccountRepository(dbConnection)
	tenantPgRepository := infra.NewPgTenantRepository(dbConnection)

	startServer(PORT, ENV, accountPgRepository, tenantPgRepository)

}

func startServer(PORT, ENV string, accountRepository port.AccountRepository, tenantRepository port.TenantRepository) {
	createAccountUsecase := accountUsecases.NewCreateAccountUsecase(accountRepository)
	findOneTenantUseCase := tenantUsecases.NewFindOneTenantUseCase(tenantRepository)

	accountHandler := handlers.NewAccountHandler(createAccountUsecase)
	tenantHandler := handlers.NewTenantHandler(findOneTenantUseCase)

	handlers := &router.Handlers{
		AccountHandler: accountHandler,
		TenantHandler:  tenantHandler,
	}

	router := router.Routes(handlers)

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))

	if err != nil {
		slog.Error("cannot start server",
			slog.String("error", err.Error()),
		)

		panic(err)
	}

	slog.Info(fmt.Sprintf("server running on port: %s - with environment: %s", PORT, ENV))
}

func initDbConnection(psqlInfo string) *sql.DB {
	slog.Info("database connection established")
	return config.InitConfig(psqlInfo)
}
