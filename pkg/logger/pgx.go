package logger

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

// safeKeys are pgx log data keys that are safe to log without redaction.
var safeKeys = map[string]bool{
	"sql":         true,
	"commandTag":  true,
	"time":        true,
	"rowCount":    true,
	"pid":         true,
	"err":         true,
	"host":        true,
	"port":        true,
	"database":    true,
	"command_tag": true,
}

// PgxLogger adapts slog to the pgx.Logger interface.
type PgxLogger struct{}

// Log implements the pgx.Logger interface using slog.
func (l *PgxLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		if safeKeys[k] {
			attrs = append(attrs, slog.Any(k, v))
		}
	}

	switch level {
	case pgx.LogLevelTrace, pgx.LogLevelDebug:
		slog.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
	case pgx.LogLevelInfo:
		slog.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
	case pgx.LogLevelWarn:
		slog.LogAttrs(ctx, slog.LevelWarn, msg, attrs...)
	case pgx.LogLevelError:
		slog.LogAttrs(ctx, slog.LevelError, msg, attrs...)
	default:
		slog.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
	}
}
