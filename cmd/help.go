package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/commander"
)

var envHelpTopics = []commander.HelpTopic{
	{
		Name:  "environment",
		Short: "List of supported environment variables",
		Long: heredoc.Doc(`
            RAYSTACK_CONFIG_DIR: the directory where stencil will store configuration files. Default:
            "$XDG_CONFIG_HOME/raystack" or "$HOME/.config/raystack".

            NO_COLOR: set to any value to avoid printing ANSI escape sequences for color output.

            CLICOLOR: set to "0" to disable printing ANSI colors in output.
        `),
	},
}
