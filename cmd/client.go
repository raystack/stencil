package cmd

import (
	"context"
	"errors"
	"time"

	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	return grpc.DialContext(ctx, host, opts...)
}

func createClient(cmd *cobra.Command) (stencilv1beta1.StencilServiceClient, func(), error) {
	host, err := cmd.Flags().GetString("host")
	if err != nil {
		return nil, nil, err
	}
	if host == "" {
		return nil, nil, errors.New("\"host\" not set")
	}

	dialTimeoutCtx, dialCancel := context.WithTimeout(cmd.Context(), time.Second*2)
	conn, err := createConnection(dialTimeoutCtx, host)
	if err != nil {
		dialCancel()
		return nil, nil, err
	}

	cancel := func() {
		dialCancel()
		conn.Close()
	}

	client := stencilv1beta1.NewStencilServiceClient(conn)
	return client, cancel, nil
}
