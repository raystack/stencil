package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"
	"github.com/spf13/cobra"
)

//New root command
func New() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "stencil <command> <subcommand> [flags]",
		Short: "Schema registry to manage schemas efficiently",
		Long: heredoc.Doc(`
			Schema registry to manage schemas efficiently.

			Stencil is a schema registry that provides schema mangement and validation to ensure data
			compatibility across applications. It enables developers to create, manage and consume 
			schemas dynamically, efficiently, and reliably, and provides a simple way to validate data 
			against those schemas. Stencil support protobuf and support for other formats coming soon.
		`),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ stencil upload 
			$ stencil download 
			$ stencil snapshot list
			$ stencil serve
			$ stencil protoc
		`),
		Annotations: map[string]string{
			"group:core": "true",
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/odpf/stencil/issues
			`),
		},
	}

	cmdx.SetHelp(cmd)

	cmd.AddCommand(ServeCmd())
	cmd.AddCommand(UploadCmd())
	cmd.AddCommand(MigrateCmd())
	cmd.AddCommand(DownloadCmd())
	cmd.AddCommand(Snapshot())
	cmd.AddCommand(ProtocCmd())
	return cmd
}
