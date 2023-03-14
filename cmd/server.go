package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/goto/stencil/config"
	"github.com/goto/stencil/internal/server"
	"github.com/goto/stencil/internal/store/postgres"
	"github.com/spf13/cobra"

	// Importing postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server <command>",
		Aliases: []string{"s"},
		Short:   "Server management",
		Long:    "Server management commands.",
		Example: heredoc.Doc(`
			$ stencil server start
			$ stencil server start -c ./config.yaml
			$ stencil server migrate
			$ stencil server migrate -c ./config.yaml
		`),
	}

	cmd.AddCommand(startCommand())
	cmd.AddCommand(migrateCommand())

	return cmd
}

func startCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"s"},
		Short:   "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}
			server.Start(cfg)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}

func migrateCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			if err := postgres.Migrate(cfg.DB.ConnectionString); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}
