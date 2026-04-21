package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/raystack/stencil/internal/store"
)

// ErrorResponse returns a new unary interceptor that standardizes error formatting.
// It ensures all errors returned from handlers are proper Connect errors.
// Non-Connect errors are sanitized to prevent leaking internal details.
func ErrorResponse() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err != nil {
				return resp, ensureConnectError(ctx, err)
			}
			return resp, nil
		}
	}
}

// ensureConnectError wraps non-Connect errors as sanitized internal Connect errors.
func ensureConnectError(ctx context.Context, err error) error {
	var connectErr *connect.Error
	if errors.As(err, &connectErr) {
		return err
	}
	// Convert known domain errors to proper Connect error codes
	var storageErr store.StorageErr
	if errors.As(err, &storageErr) {
		return storageErr.ConnectError()
	}
	ref := time.Now().Unix()
	slog.ErrorContext(ctx, "unhandled error", "error", err, "ref", ref)
	return connect.NewError(connect.CodeInternal, fmt.Errorf(
		"%s - ref (%d)",
		http.StatusText(http.StatusInternalServerError),
		ref,
	))
}
