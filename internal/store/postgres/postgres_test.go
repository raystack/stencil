package postgres_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/raystack/stencil/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func tearDown(t *testing.T) {
	t.Helper()
	connectionString := os.Getenv("TEST_DB_CONNECTIONSTRING")
	if connectionString == "" {
		t.Skip("Skipping test since DB info not available")
		return
	}
	m, err := postgres.NewHTTPFSMigrator(connectionString)
	if assert.NoError(t, err) {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			t.Fatalf("failed to rollback migrations: %v", err)
		}
	}
}
