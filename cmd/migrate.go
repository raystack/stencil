package cmd

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/odpf/stencil/server/config"

	// Importing postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		Run:   migrateCmd,
	})
}

func migrateCmd(cmd *cobra.Command, args []string) {
	c := config.LoadConfig()
	m, err := migrate.New(
		c.DB.MigrationsPath,
		c.DB.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
