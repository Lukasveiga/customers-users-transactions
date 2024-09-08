package factory

import (
	"database/sql"
	"fmt"

	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/handlers"
	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository/sqlc"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	cardUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/cards"
	tenantUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
)

type Handlers struct {
	AccountHandler *handlers.AccountHandler
	TenantHandler  *handlers.TenantHandler
	CardHandler    *handlers.CardHandler
}

func InitHandlers(dbConnection *sql.DB) *Handlers {
	// Repository
	repository := infra.NewTx(dbConnection)

	// Account usecases
	createAccountUsecase := accountUsecases.NewCreateAccountUsecase(repository)
	findOneAccountUsecase := accountUsecases.NewFindOneAccountUsecase(repository)
	findAllAccountsUsecase := accountUsecases.NewFindAllAccountsUsecase(repository)
	updateAccountUsecase := accountUsecases.NewActiveAccountUsecase(repository)
	deleteAaccountUsecase := accountUsecases.NewInactiveAccountUsecase(repository)

	// Tenant usecases
	findOneTenantUseCase := tenantUsecases.NewFindOneTenantUseCase(repository)

	// Card usecases
	createCardUsecase := cardUsecases.NewCreateCardUsecase(repository, findOneAccountUsecase)
	findCardUsecase := cardUsecases.NewFindCardUsecase(repository, findOneAccountUsecase)

	// Handlers
	accountHandler := handlers.NewAccountHandler(createAccountUsecase, findAllAccountsUsecase, findOneAccountUsecase, updateAccountUsecase, deleteAaccountUsecase)
	tenantHandler := handlers.NewTenantHandler(findOneTenantUseCase)
	cardHandler := handlers.NewCardHandler(createCardUsecase, findCardUsecase)

	return &Handlers{
		AccountHandler: accountHandler,
		TenantHandler:  tenantHandler,
		CardHandler:    cardHandler,
	}
}

func GetDbUrlConn(env string) string {
	var (
		host, db_port, user, password, dbname string
	)

	switch env {
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

	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, db_port, user, password, dbname)
}
