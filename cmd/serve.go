package cmd

import (
	"github.com/odpf/stencil/config"
	"github.com/odpf/stencil/server"
	"github.com/spf13/cobra"
)

// ServeCmd start new stencil server
func ServeCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Run stencil server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}
			server.Start(cfg)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}
