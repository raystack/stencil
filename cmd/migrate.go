package cmd

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/odpf/salt/config"
	"github.com/odpf/stencil/server"

	// Importing postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

// MigrateCmd start new stencil server
func MigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg server.Config
			loader := config.NewLoader(config.WithPath("./"))
			if err := loader.Load(&cfg); err != nil {
				log.Fatal(err)
			}
			m, err := migrate.New(
				cfg.DB.MigrationsPath,
				cfg.DB.ConnectionString)
			if err != nil {
				log.Fatal(err)
			}
			if err := m.Up(); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
}
