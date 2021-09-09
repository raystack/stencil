package cmd

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var cmd = &cobra.Command{
		Use:           "stencil <command> <subcommand> [flags]",
		Short:         "Dynamic schema registry",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(ServeCmd())
	cmd.AddCommand(UploadCmd())
	cmd.AddCommand(MigrateCmd())
	cmd.AddCommand(DownloadCmd())
	return cmd
}
