package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/commander"
	"github.com/raystack/salt/config"
	"github.com/spf13/cobra"
)

type CDK struct {
	Config *config.Loader
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
				Read the manual at https://raystack.github.io/stencil/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/raystack/stencil/issues
			`),
		},
	}

	cdk := &CDK{
		Config: config.NewLoader(
			config.WithAppConfig("stencil"),
			config.WithEnvPrefix("STENCIL"),
			config.WithFlags(cmd.Flags()),
		),
	}

	cmd.AddCommand(ServerCommand())
	cmd.AddCommand(configCmd(cdk))
	cmd.AddCommand(NamespaceCmd(cdk))
	cmd.AddCommand(SchemaCmd(cdk))
	cmd.AddCommand(SearchCmd(cdk))

	hooks := []commander.HookBehavior{
		{
			Name: "client",
			Behavior: func(cmd *cobra.Command) {
				cmd.PersistentFlags().String("host", "", "Server host address")
			},
		},
	}

	// Help topics
	cmdr := commander.New(
		cmd,
		commander.WithTopics(envHelpTopics),
		commander.WithHooks(hooks),
	)
	cmdr.Init()

	return cmd
}
