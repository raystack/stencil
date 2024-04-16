package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goto/stencil/core/namespace"
	"github.com/goto/stencil/core/schema"
	"github.com/goto/stencil/internal/store/postgres"
)

func getSchemaStore(t *testing.T) *postgres.SchemaRepository {
	t.Helper()
	connectionString := os.Getenv("TEST_DB_CONNECTIONSTRING")
	if connectionString == "" {
		t.Skip("Skipping test since DB info not available")
		return nil
	}
	err := postgres.Migrate(connectionString)
	assert.Nil(t, err)
	dbc := postgres.NewStore(connectionString)
	return postgres.NewSchemaRepository(dbc)
}

func TestSchema(t *testing.T) {
	tearDown(t)
	db := getSchemaStore(t)
	namespaceStore := getNamespaceStore(t)
	ctx := context.Background()

	t.Run("schema", func(t *testing.T) {
		n := &namespace.Namespace{ID: "testschema", Format: "protobuf", Compatibility: "FULL", Description: "testDesc"}
		_, err := namespaceStore.Create(ctx, *n)
		assert.Nil(t, err)
		meta := &schema.Metadata{
			Format: "avro",
		}
		t.Run("create: should create schema", func(t *testing.T) {
			versionNumber, err := db.Create(ctx, n.ID, "sName", meta, "uuid-1", &schema.SchemaFile{ID: "t1", Data: []byte("testdata")})
			assert.Nil(t, err)
			assert.Equal(t, int32(1), versionNumber)
		})
		t.Run("create: should increment version number on new schema", func(t *testing.T) {
			versionNumber, err := db.Create(ctx, n.ID, "sName", meta, "uuid-2", &schema.SchemaFile{ID: "t2", Data: []byte("testdata-2")})
			assert.Nil(t, err)
			assert.Equal(t, int32(2), versionNumber)
		})
		t.Run("create: should return same version number if schema is same", func(t *testing.T) {
			versionNumber, err := db.Create(ctx, n.ID, "sName", meta, "uuid-1", &schema.SchemaFile{ID: "t1", Data: []byte("testdata")})
			assert.Nil(t, err)
			assert.Equal(t, int32(1), versionNumber)
		})
		t.Run("list_schemas: should return schema", func(t *testing.T) {
			schemaList, err := db.List(ctx, "testschema")
			assert.Nil(t, err)
			assert.Equal(t, []schema.Schema{{Name: "sName", Format: "avro", Compatibility: "", Authority: ""}}, schemaList)
		})
		t.Run("list_versions: should return versions for specified schema", func(t *testing.T) {
			schemaList, err := db.ListVersions(ctx, "testschema", "sName")
			assert.Nil(t, err)
			assert.Equal(t, []int32{1, 2}, schemaList)
		})
		t.Run("get: should return specified schema", func(t *testing.T) {
			s, err := db.Get(ctx, n.ID, "sName", 1)
			assert.Nil(t, err)
			assert.Equal(t, []byte("testdata"), s)
		})
		t.Run("get: should return specified schema ID", func(t *testing.T) {
			s, err := db.GetSchemaID(ctx, n.ID, "sName")
			assert.Nil(t, err)
			assert.NotZero(t, s)
		})
		t.Run("getMetadata: should return metadata", func(t *testing.T) {
			actual, err := db.GetMetadata(ctx, n.ID, "sName")
			assert.Nil(t, err)
			assert.Equal(t, meta.Format, actual.Format)
		})
		t.Run("updateMetadata: should update metadata", func(t *testing.T) {
			actual, err := db.UpdateMetadata(ctx, n.ID, "sName", &schema.Metadata{Compatibility: "FULL"})
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
			err := db.Delete(ctx, n.ID, "sName")
			assert.Nil(t, err)
			schemaList, err := db.List(ctx, "testschema")
			assert.Nil(t, err)
			assert.Equal(t, 0, len(schemaList))
		})
	})
	tearDown(t)
}
