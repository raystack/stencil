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
		Annotations: map[string]string{
			"group": "core",
			"help:learn": heredoc.Doc(`
				Use 'stencil <command> --help' for info about a command.
				Read the manual at https://odpf.github.io/stencil/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/odpf/stencil/issues
			`),
		},
	}

	cmd.AddCommand(ServerCommand())
	cmd.AddCommand(NamespaceCmd())
	cmd.AddCommand(SchemaCmd())
	cmd.AddCommand(SearchCmd())

	// Help topics
	cmdx.SetHelp(cmd)
	cmd.AddCommand(cmdx.SetCompletionCmd("stencil"))
	cmd.AddCommand(cmdx.SetHelpTopicCmd("environment", envHelp))
	cmd.AddCommand(cmdx.SetRefCmd(cmd))

	return cmd
}
