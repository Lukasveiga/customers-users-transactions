package config

import (
	"io"
	"log/slog"
	"os"
)

type LoggerConfig struct {
	Env     string
	LogPath string
}

func InitLogger(loggerConfig *LoggerConfig) {
	level := new(slog.LevelVar)
	var writer io.Writer

	switch loggerConfig.Env {
	case "prod":
		level.Set(slog.LevelInfo.Level())
		writer = os.Stdout
	default:
		level.Set(slog.LevelDebug.Level())
		writer = os.Stdout
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	l := slog.New(slog.NewJSONHandler(writer, handlerOpts))
	slog.SetDefault(l)
}
