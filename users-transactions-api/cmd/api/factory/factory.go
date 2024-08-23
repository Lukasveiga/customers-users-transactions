package factory

import (
	"database/sql"
	"fmt"

	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/Lukasveiga/customers-users-transaction/internal/handlers"
	infra "github.com/Lukasveiga/customers-users-transaction/internal/infra/repository"
	accountUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/account"
	tenantUsecases "github.com/Lukasveiga/customers-users-transaction/internal/usecases/tenant"
)

type Handlers struct {
	AccountHandler *handlers.AccountHandler
	TenantHandler  *handlers.TenantHandler
}

func InitHandlers(dbConnection *sql.DB) *Handlers {
	accountPgRepository := infra.NewPgAccountRepository(dbConnection)
	tenantPgRepository := infra.NewPgTenantRepository(dbConnection)

	createAccountUsecase := accountUsecases.NewCreateAccountUsecase(accountPgRepository)
	findOneAccountUsecase := accountUsecases.NewFindOneAccountUsecase(accountPgRepository)
	findAllAccountsUsecase := accountUsecases.NewFindAllAccountsUsecase(accountPgRepository)
	updateAccountUsecase := accountUsecases.NewUpdateAccountUsecase(accountPgRepository)
	deleteAaccountUsecase := accountUsecases.NewDeleteAccountUsecase(accountPgRepository)

	findOneTenantUseCase := tenantUsecases.NewFindOneTenantUseCase(tenantPgRepository)

	accountHandler := handlers.NewAccountHandler(createAccountUsecase, findAllAccountsUsecase, findOneAccountUsecase, updateAccountUsecase, deleteAaccountUsecase)
	tenantHandler := handlers.NewTenantHandler(findOneTenantUseCase)

	return &Handlers{
		AccountHandler: accountHandler,
		TenantHandler:  tenantHandler,
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
