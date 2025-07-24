package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	LogLevel   slog.Leveler
	JSONFormat bool
}

func Init(cfg Config) {
	var handler slog.Handler

	options := &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.LogLevel,
	}

	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
