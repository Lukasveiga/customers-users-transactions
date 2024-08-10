package config

import (
	"database/sql"
	"log/slog"
)

func InitConfig(connString string) *sql.DB {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		slog.Error(
			"database configuration",
			slog.String("error", err.Error()),
		)

		panic(err)
	}

	return db
}
