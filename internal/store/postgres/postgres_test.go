package postgres_test

import (
	"os"
	"testing"

	"github.com/goto/stencil/internal/store/postgres"
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
		m.Down()
	}
}
