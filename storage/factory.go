package storage

import (
	"fmt"
	"strings"

	"github.com/odpf/stencil/config"
	"github.com/odpf/stencil/storage/postgres"
	"github.com/odpf/stencil/storage/sqlite"
)

// FactoryStore creates store based on config
func FactoryStore(cfg config.DBConfig) Store {
	dbType := strings.ToLower(cfg.Type)
	if dbType == "postgres" {
		return postgres.NewStore(cfg.ConnectionString)
	}
	if dbType == "sqlite" {
		connectionString := "file:" + cfg.ConnectionString
		return sqlite.NewStore(connectionString)
	}
	panic(fmt.Sprintf("database with type [%s] is not recognized", dbType))
}

// FactoryMigrate migrates to store based on config
func FactoryMigrate(cfg config.DBConfig) error {
	dbType := strings.ToLower(cfg.Type)
	if dbType == "postgres" {
		return postgres.Migrate(cfg.ConnectionString)
	}
	if dbType == "sqlite" {
		connectionString := "sqlite3://" + cfg.ConnectionString
		return sqlite.Migrate(connectionString)
	}
	panic(fmt.Sprintf("database with type [%s] is not recognized", dbType))
}
