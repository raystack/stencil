package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"
	"github.com/spf13/cobra"
)

type Config struct {
	Host string `mapstructure:"host"`
}

func LoadConfig() (*Config, error) {
	var config Config

	cfg := cmdx.SetConfig("stencil")
	err := cfg.Load(&config)

	return &config, err
}

func configCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage client configurations",
		Example: heredoc.Doc(`
			$ stencil config init
			$ stencil config list`),
	}

	cmd.AddCommand(configInitCommand())
	cmd.AddCommand(configListCommand())

	return cmd
}

func configInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new client configuration",
		Example: heredoc.Doc(`
			$ stencil config init
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmdx.SetConfig("stencil")

			if err := cfg.Init(&Config{}); err != nil {
				return err
			}

			fmt.Printf("config created: %v\n", cfg.File())
			return nil
		},
	}
}

func configListCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List client configuration settings",
		Example: heredoc.Doc(`
			$ stencil config list
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cmdx.SetConfig("stencil")

			data, err := cfg.Read()
			if err != nil {
				return ErrClientConfigNotFound
			}

			fmt.Println(data)
			return nil
		},
	}
	return cmd
}
