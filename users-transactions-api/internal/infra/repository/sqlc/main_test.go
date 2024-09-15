package infra

import (
	"database/sql"
	"os"
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/config"
)

var (
	testQueries *Queries
	testDb      *sql.DB
)

func TestMain(m *testing.M) {
	testDb = config.SetupPgTestcontainers()

	testQueries = New(testDb)

	os.Exit(m.Run())
}
