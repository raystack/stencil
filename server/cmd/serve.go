package cmd

import (
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run:   serve,
	})
}

func serve(cmd *cobra.Command, args []string) {
	c := config.LoadConfig()
	server.Start(c)
}
