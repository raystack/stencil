package cmd

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

var (
	ErrClientConfigNotFound = errors.New(heredoc.Doc(`
		Stencil client config not found.

		Run "stencil config init" to initialize a new client config or
		Run "stencil help environment" for more information.
	`))
	ErrClientConfigHostNotFound = errors.New(heredoc.Doc(`
		Stencil client config "host" not found.

		Pass stencil server host with "--host" flag or 
		set host in stencil config.

		Run "stencil config <subcommand>" or
		"stencil help environment" for more information.
	`))
	ErrClientNotAuthorized = errors.New(heredoc.Doc(`
		Stencil auth error. Stencil requires an auth header.
		
		Run "stencil help auth" for more information.
	`))
)
