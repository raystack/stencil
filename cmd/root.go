package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/cmdx"
	"github.com/spf13/cobra"
)

type CDK struct {
	Config *cmdx.Config
}

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
				Read the manual at https://goto.github.io/stencil/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/goto/stencil/issues
			`),
		},
	}

	cdk := &CDK{Config: cmdx.SetConfig("stencil")}

	cmd.AddCommand(ServerCommand())
	cmd.AddCommand(configCmd(cdk))
	cmd.AddCommand(NamespaceCmd(cdk))
	cmd.AddCommand(SchemaCmd(cdk))
	cmd.AddCommand(SearchCmd(cdk))

	// Help topics
	cmdx.SetHelp(cmd)
	cmd.AddCommand(cmdx.SetCompletionCmd("stencil"))
	cmd.AddCommand(cmdx.SetHelpTopicCmd("environment", envHelp))
	cmd.AddCommand(cmdx.SetRefCmd(cmd))

	cmdx.SetClientHook(cmd, func(cmd *cobra.Command) {
		// client config
		cmd.PersistentFlags().String("host", "", "Server host address")
	})

	return cmd
}
