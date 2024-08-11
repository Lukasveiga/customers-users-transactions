package config

import (
	"context"
	"database/sql"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	ctx        context.Context
	connString string
)

func setupContainer() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	ctx = context.Background()
	c, err := postgres.Run(
		ctx,
		"postgres:14-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgre"),
		postgres.WithPassword("postgre"),
		postgres.WithInitScripts(basepath+"/init_db.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		slog.Error("testcontainers configuration - RunContainer", "error", err)
		panic(err)
	}

	connString, err = c.ConnectionString(ctx)

	connString = connString + "sslmode=disable"

	slog.Debug("testcontainers url connection", "url", connString)

	if err != nil {
		slog.Error("testcontainers configuration - ConnectionString", "error", err)
		panic(err)
	}
}

func SetupPgTestcontainers() *sql.DB {
	setupContainer()
	return InitConfig(connString)
}
