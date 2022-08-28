package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/odpf/stencil/core/namespace"
	"github.com/odpf/stencil/internal/store"
	"github.com/odpf/stencil/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func getNamespaceStore(t *testing.T) *postgres.NamespaceRepository {
	t.Helper()
	connectionString := os.Getenv("TEST_DB_CONNECTIONSTRING")
	if connectionString == "" {
		t.Skip("Skipping test since DB info not available")
		return nil
	}
	err := postgres.Migrate(connectionString)
	assert.Nil(t, err)
	dbc := postgres.NewStore(connectionString)
	return postgres.NewNamespaceRepository(dbc)
}

func TestNamespace(t *testing.T) {
	tearDown(t)
	db := getNamespaceStore(t)
	ctx := context.Background()
	n := &namespace.Namespace{ID: "test", Format: "protobuf", Compatibility: "FULL", Description: "testDesc"}
	t.Run("Namespace", func(t *testing.T) {
		t.Run("create: should create namespace", func(t *testing.T) {
			ns, err := db.Create(ctx, *n)
			assert.Nil(t, err)
			assertNamespace(t, *n, ns)
		})
		t.Run("create: should return error on duplicate namespace name", func(t *testing.T) {
			_, err := db.Create(ctx, *n)
			assert.ErrorIs(t, err, store.ConflictErr)
		})
		t.Run("update: should update the namespace", func(t *testing.T) {
			n.Description = "newDescription"
			n.Format = "avro"
			ns, err := db.Update(ctx, *n)
			assert.Nil(t, err)
			assertNamespace(t, *n, ns)
		})
		t.Run("update: should return error if namespace not found", func(t *testing.T) {
			n.ID = "test2"
			_, err := db.Update(ctx, *n)
			assert.ErrorIs(t, err, store.NoRowsErr)
			n.ID = "test"
		})
		t.Run("get: should get the namespace", func(t *testing.T) {
			ns, err := db.Get(ctx, "test")
			assert.Nil(t, err)
			assertNamespace(t, *n, ns)
		})
		t.Run("get: should return the error if namespace not found", func(t *testing.T) {
			_, err := db.Get(ctx, "test1")
			assert.ErrorIs(t, err, store.NoRowsErr)
		})
		t.Run("list: should list created namespaces", func(t *testing.T) {
			ls, err := db.List(ctx)
			assert.Nil(t, err)
			assert.Equal(t, []string{"test"}, ls)
		})
		t.Run("delete: should delete namespace", func(t *testing.T) {
			err := db.Delete(ctx, "test")
			assert.Nil(t, err)
		})
	})

}

func assertNamespace(t *testing.T, expected, actual namespace.Namespace) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Compatibility, actual.Compatibility)
	assert.Equal(t, expected.Format, actual.Format)
	assert.Equal(t, expected.Description, actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}
