package middleware

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

// Logger returns a new unary interceptor that logs request details.
func Logger() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			resp, err := next(ctx, req)

			duration := time.Since(start)
			attrs := []slog.Attr{
				slog.String("procedure", req.Spec().Procedure),
				slog.Duration("duration", duration),
				slog.String("peer_addr", req.Peer().Addr),
			}

			if err != nil {
				attrs = append(attrs, slog.Any("error", err))
				connectErr, ok := err.(*connect.Error)
				if ok {
					attrs = append(attrs, slog.String("code", connectErr.Code().String()))
				}
				slog.LogAttrs(ctx, slog.LevelError, "request failed", attrs...)
			} else {
				slog.LogAttrs(ctx, slog.LevelDebug, "request completed", attrs...)
			}

			return resp, err
		}
	}
}
