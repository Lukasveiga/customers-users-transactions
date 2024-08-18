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
		file, err := os.OpenFile(loggerConfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		level.Set(slog.LevelInfo.Level())
		writer = io.MultiWriter(os.Stdout, file)
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
