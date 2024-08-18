package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

func GetEnv(key string) string {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		slog.Error("environment configuration",
			slog.String("error", err.Error()))
		panic(err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		slog.Error("environment configuration",
			slog.String("error", "Invalid type assertion"))
		panic(err)
	}

	return value
}
