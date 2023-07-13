package postgres

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/raystack/stencil/internal/store"
	"github.com/raystack/stencil/pkg/logger"
)

//go:embed migrations
var migrationFs embed.FS

const (
	resourcePath = "migrations"
)

// DB represents postgres database instance
type DB struct {
	*pgxpool.Pool
}

// NewStore create a postgres store
func NewStore(conn string) *DB {
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

func wrapError(err error, format string, args ...interface{}) error {
	if err == nil {
		return err
	}
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		return store.NoRowsErr.WithErr(err, fmt.Sprintf(format, args...))
	}
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return store.ConflictErr.WithErr(err, fmt.Sprintf(format, args...))
		}
	}
	return store.UnknownErr.WithErr(err, fmt.Sprintf(format, args...))
}
