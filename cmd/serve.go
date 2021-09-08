package cmd

import (
	"github.com/odpf/stencil/server"
	"github.com/spf13/cobra"
)

// ServeCmd start new stencil server
func ServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Aliases: []string{"v"},
		Short:   "Run stencil server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server.Start()
			return nil
		},
	}
}
