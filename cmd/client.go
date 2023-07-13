package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/raystack/salt/cmdx"
	"github.com/raystack/salt/config"
	stencilv1beta1 "github.com/raystack/stencil/proto/raystack/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Host string `yaml:"host" cmdx:"host"`
}

func createConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	return grpc.DialContext(ctx, host, opts...)
}

func createClient(cmd *cobra.Command, cdk *CDK) (stencilv1beta1.StencilServiceClient, func(), error) {
	c, err := loadClientConfig(cmd, cdk.Config)
	if err != nil {
		return nil, nil, err
	}

	host := c.Host

	if host == "" {
		return nil, nil, ErrClientConfigHostNotFound
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

func loadClientConfig(cmd *cobra.Command, cmdxConfig *cmdx.Config) (*ClientConfig, error) {
	var clientConfig ClientConfig

	if err := cmdxConfig.Load(
		&clientConfig,
		cmdx.WithFlags(cmd.Flags()),
	); err != nil {
		if !errors.Is(err, new(config.ConfigFileNotFoundError)) {
			return nil, ErrClientConfigNotFound
		}
	}

	return &clientConfig, nil
}
