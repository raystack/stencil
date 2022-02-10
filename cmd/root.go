package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"
	"github.com/spf13/cobra"
)

// New root command
func New() *cobra.Command {
	var cmd = &cobra.Command{
		Use:           "stencil <command> <subcommand> [flags]",
		Short:         "Schema registry",
		Long:          "Schema registry to manage schemas efficiently.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ stencil namespace create
			$ stencil schema create
			$ stencil upload
			$ stencil download
			$ stencil snapshot list
			$ stencil server start
		`),
		Annotations: map[string]string{
			"group:core": "true",
			"help:learn": heredoc.Doc(`
				Use 'stencil <command> <subcommand> --help' for more information about a command.
				Read the manual at https://odpf.github.io/stencil/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/odpf/stencil/issues
			`),
		},
	}

	cmd.AddCommand(ServerCommand())
	cmd.AddCommand(GraphCmd())
	cmd.AddCommand(PrintCmd())
	cmd.AddCommand(NamespaceCmd())
	cmd.AddCommand(SchemaCmd())

	// Help topics
	cmdx.SetHelp(cmd)
	cmd.AddCommand(cmdx.SetCompletionCmd("stencil"))
	cmd.AddCommand(cmdx.SetHelpTopic("environment", envHelp))
	cmd.AddCommand(cmdx.SetRefCmd(cmd))

	return cmd
}
