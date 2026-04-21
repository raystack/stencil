package middleware

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"connectrpc.com/connect"
)

// Recovery returns a new unary interceptor that recovers from panics.
func Recovery() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					log.Printf("panic recovered: %v\n%s", r, string(stack))
					err = connect.NewError(
						connect.CodeInternal,
						fmt.Errorf("internal server error"),
					)
				}
			}()
			return next(ctx, req)
		}
	}
}
