package cmd

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/odpf/stencil/config"

	// Importing postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

// MigrateCmd start new stencil server
func MigrateCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
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

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}
