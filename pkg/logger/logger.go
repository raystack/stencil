package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init sets up the global slog logger with a JSON handler.
func Init(logLevel string) {
	var level slog.LevelVar
	switch strings.ToLower(logLevel) {
	case "debug":
		level.Set(slog.LevelDebug)
	case "warn", "warning":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	default:
		level.Set(slog.LevelInfo)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     &level,
		AddSource: true,
	}))
	slog.SetDefault(logger)
}
