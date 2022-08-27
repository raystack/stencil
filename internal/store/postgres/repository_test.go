package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/odpf/stencil/core/namespace"
	"github.com/odpf/stencil/domain"
	"github.com/odpf/stencil/internal/store"
	"github.com/odpf/stencil/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func getStore(t *testing.T) *postgres.Store {
	t.Helper()
	connectionString := os.Getenv("TEST_DB_CONNECTIONSTRING")
	if connectionString == "" {
		t.Skip("Skipping test since DB info not available")
		return nil
	}
	err := postgres.Migrate(connectionString)
	assert.Nil(t, err)
	return postgres.NewStore(connectionString)
}

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

func assertNamespace(t *testing.T, expected, actual namespace.Namespace) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Compatibility, actual.Compatibility)
	assert.Equal(t, expected.Format, actual.Format)
	assert.Equal(t, expected.Description, actual.Description)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}

func TestStorage(t *testing.T) {
	tearDown(t)
	db := getStore(t)
	ctx := context.Background()
	n := &namespace.Namespace{ID: "test", Format: "protobuf", Compatibility: "FULL", Description: "testDesc"}
	t.Run("Namespace", func(t *testing.T) {
		t.Run("create: should create namespace", func(t *testing.T) {
			ns, err := db.CreateNamespace(ctx, *n)
			assert.Nil(t, err)
			assertNamespace(t, *n, ns)
		})
		t.Run("create: should return error on duplicate namespace name", func(t *testing.T) {
			_, err := db.CreateNamespace(ctx, *n)
			assert.ErrorIs(t, err, store.ConflictErr)
		})
		t.Run("update: should update the namespace", func(t *testing.T) {
			n.Description = "newDescription"
			n.Format = "avro"
			ns, err := db.UpdateNamespace(ctx, *n)
			assert.Nil(t, err)
			assertNamespace(t, *n, ns)
		})
		t.Run("update: should return error if namespace not found", func(t *testing.T) {
			n.ID = "test2"
			_, err := db.UpdateNamespace(ctx, *n)
			assert.ErrorIs(t, err, store.NoRowsErr)
			n.ID = "test"
		})
		t.Run("get: should get the namespace", func(t *testing.T) {
			ns, err := db.GetNamespace(ctx, "test")
			assert.Nil(t, err)
			assertNamespace(t, *n, ns)
		})
		t.Run("get: should return the error if namespace not found", func(t *testing.T) {
			_, err := db.GetNamespace(ctx, "test1")
			assert.ErrorIs(t, err, store.NoRowsErr)
		})
		t.Run("list: should list created namespaces", func(t *testing.T) {
			ls, err := db.ListNamespaces(ctx)
			assert.Nil(t, err)
			assert.Equal(t, []string{"test"}, ls)
		})
		t.Run("delete: should delete namespace", func(t *testing.T) {
			err := db.DeleteNamespace(ctx, "test")
			assert.Nil(t, err)
		})
	})

	t.Run("schema", func(t *testing.T) {
		n := &namespace.Namespace{ID: "testschema", Format: "protobuf", Compatibility: "FULL", Description: "testDesc"}
		_, err := db.CreateNamespace(ctx, *n)
		assert.Nil(t, err)
		meta := &domain.Metadata{
			Format: "avro",
		}
		t.Run("create: should create schema", func(t *testing.T) {
			versionNumber, err := db.CreateSchema(ctx, n.ID, "sName", meta, "uuid-1", &domain.SchemaFile{ID: "t1", Data: []byte("testdata")})
			assert.Nil(t, err)
			assert.Equal(t, int32(1), versionNumber)
		})
		t.Run("create: should increment version number on new schema", func(t *testing.T) {
			versionNumber, err := db.CreateSchema(ctx, n.ID, "sName", meta, "uuid-2", &domain.SchemaFile{ID: "t2", Data: []byte("testdata-2")})
			assert.Nil(t, err)
			assert.Equal(t, int32(2), versionNumber)
		})
		t.Run("create: should return same version number if schema is same", func(t *testing.T) {
			versionNumber, err := db.CreateSchema(ctx, n.ID, "sName", meta, "uuid-1", &domain.SchemaFile{ID: "t1", Data: []byte("testdata")})
			assert.Nil(t, err)
			assert.Equal(t, int32(1), versionNumber)
		})
		t.Run("list_schemas: should return schema", func(t *testing.T) {
			schemaList, err := db.ListSchemas(ctx, "testschema")
			assert.Nil(t, err)
			assert.Equal(t, []string{"sName"}, schemaList)
		})
		t.Run("list_versions: should return versions for specified schema", func(t *testing.T) {
			schemaList, err := db.ListVersions(ctx, "testschema", "sName")
			assert.Nil(t, err)
			assert.Equal(t, []int32{1, 2}, schemaList)
		})
		t.Run("get: should return specified schema", func(t *testing.T) {
			s, err := db.GetSchema(ctx, n.ID, "sName", 1)
			assert.Nil(t, err)
			assert.Equal(t, []byte("testdata"), s)
		})
		t.Run("getMetadata: should return metadata", func(t *testing.T) {
			actual, err := db.GetSchemaMetadata(ctx, n.ID, "sName")
			assert.Nil(t, err)
			assert.Equal(t, meta.Format, actual.Format)
		})
		t.Run("updateMetadata: should update metadata", func(t *testing.T) {
			actual, err := db.UpdateSchemaMetadata(ctx, n.ID, "sName", &domain.Metadata{Compatibility: "FULL"})
			assert.Nil(t, err)
			assert.Equal(t, "FULL", actual.Compatibility)
		})
		t.Run("getLatestVersion: should return latest schema version", func(t *testing.T) {
			s, err := db.GetLatestVersion(ctx, n.ID, "sName")
			assert.Nil(t, err)
			assert.Equal(t, int32(2), s)
		})
		t.Run("deleteVersion: should delete specified version schema", func(t *testing.T) {
			err := db.DeleteVersion(ctx, n.ID, "sName", int32(2))
			assert.Nil(t, err)
			schemaList, err := db.ListVersions(ctx, "testschema", "sName")
			assert.Nil(t, err)
			assert.Equal(t, []int32{1}, schemaList)
		})

		t.Run("deleteSchema: should delete specified schema", func(t *testing.T) {
			err := db.DeleteSchema(ctx, n.ID, "sName")
			assert.Nil(t, err)
			schemaList, err := db.ListSchemas(ctx, "testschema")
			assert.Nil(t, err)
			assert.Equal(t, 0, len(schemaList))
		})
	})
	tearDown(t)
}
