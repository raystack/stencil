package postgres_test

import (
	"os"
	"testing"

	"github.com/odpf/stencil/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func getSearchStore(t *testing.T) *postgres.SearchRepository {
	t.Helper()
	connectionString := os.Getenv("TEST_DB_CONNECTIONSTRING")
	if connectionString == "" {
		t.Skip("Skipping test since DB info not available")
		return nil
	}
	err := postgres.Migrate(connectionString)
	assert.Nil(t, err)
	dbc := postgres.NewStore(connectionString)
	return postgres.NewSearchRepository(dbc)
}
