package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

//go:embed migrations
var migrationFs embed.FS

const (
	resourcePath = "migrations"
)

// NewStore initializes sqlite store
func NewStore(connection string) *Store {
	db, err := sql.Open("sqlite3", connection)
	if err != nil {
		log.Fatal(err)
	}
	return &Store{
		db: db,
	}
}

// NewHTTPFSMigrator reads the migrations from httpfs and returns the migrate.Migrate
func NewHTTPFSMigrator(DBConnURL string) (*migrate.Migrate, error) {
	src, err := httpfs.New(http.FS(migrationFs), resourcePath)
	if err != nil {
		return &migrate.Migrate{}, fmt.Errorf("db migrator: %v", err)
	}
	return migrate.NewWithSourceInstance("httpfs", src, DBConnURL)
}

// Migrate migrates to sqlite store
func Migrate(connURL string) error {
	m, err := NewHTTPFSMigrator(connURL)
	if err != nil {
		return errors.Wrap(err, "db migrator")
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "db migrator")
	}
	return nil
}
