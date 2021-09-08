package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/odpf/salt/config"
	"github.com/odpf/stencil/server"

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
	var cfg server.Config
	loader := config.NewLoader(config.WithPath("./"))

	if err := loader.Load(&cfg); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
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
}
