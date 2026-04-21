package cmd

import (
	"net/http"

	"github.com/raystack/salt/config"
	stencilv1beta1connect "github.com/raystack/stencil/gen/raystack/stencil/v1beta1/stencilv1beta1connect"
	"github.com/spf13/cobra"
)

type ClientConfig struct {
	Host string `yaml:"host" cmdx:"host"`
}

func createClient(cmd *cobra.Command, cdk *CDK) (stencilv1beta1connect.StencilServiceClient, error) {
	c, err := loadClientConfig(cmd, cdk.Config)
	if err != nil {
		return nil, err
	}

	host := c.Host

	if host == "" {
		return nil, ErrClientConfigHostNotFound
	}

	client := stencilv1beta1connect.NewStencilServiceClient(
		http.DefaultClient,
		host,
	)
	return client, nil
}

func loadClientConfig(cmd *cobra.Command, cmdxConfig *config.Loader) (*ClientConfig, error) {
	var clientConfig ClientConfig

	if err := cmdxConfig.Load(
		&clientConfig,
	); err != nil {
		return nil, err
	}

	return &clientConfig, nil
}
