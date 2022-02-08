package cmd

import (
	"github.com/MakeNowJust/heredoc"
)

var envHelp = map[string]string{
	"short": "Environment variables that can be used with guardian",
	"long": heredoc.Doc(`
			ODPF_CONFIG_DIR: the directory where guardian will store configuration files. Default:
			"$XDG_CONFIG_HOME/odpf" or "$HOME/.config/odpf".

			NO_COLOR: set to any value to avoid printing ANSI escape sequences for color output.

			CLICOLOR: set to "0" to disable printing ANSI colors in output.
		`),
}
