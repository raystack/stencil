package store

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/odpf/stencil/server/logger"
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
