package store

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/odpf/stencil/server/config"
	"go.uber.org/zap"
)

// DB db instance
type DB struct {
	*pgxpool.Pool
}

// NewDBStore create DB store
func NewDBStore(dbConfig *config.Config) *DB {
	cc, _ := pgxpool.ParseConfig(dbConfig.DB.ConnectionString)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	cc.ConnConfig.Logger = zapadapter.NewLogger(logger)
	cc.ConnConfig.LogLevel = pgx.LogLevelInfo

	pgxPool, err := pgxpool.ConnectConfig(context.Background(), cc)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{Pool: pgxPool}
}
