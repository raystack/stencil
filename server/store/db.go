package store

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/odpf/stencil/server/logger"
	"github.com/pkg/errors"
)

//go:embed migrations
var migrationFs embed.FS

const (
	resourcePath = "migrations"
)

// DB db instance
type DB struct {
	*pgxpool.Pool
}

// NewDBStore create DB store
func NewDBStore(conn string) *DB {
	cc, _ := pgxpool.ParseConfig(conn)
	cc.ConnConfig.Logger = zapadapter.NewLogger(logger.Logger)

	pgxPool, err := pgxpool.ConnectConfig(context.Background(), cc)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{Pool: pgxPool}
}

// NewHTTPFSMigrator reads the migrations from httpfs and returns the migrate.Migrate
func NewHTTPFSMigrator(DBConnURL string) (*migrate.Migrate, error) {
	src, err := httpfs.New(http.FS(migrationFs), resourcePath)
	if err != nil {
		return &migrate.Migrate{}, fmt.Errorf("db migrator: %v", err)
	}
	return migrate.NewWithSourceInstance("httpfs", src, DBConnURL)
}

// Migrate to run up migrations
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
