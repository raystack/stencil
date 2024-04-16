package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goto/stencil/core/changedetector"
	"github.com/goto/stencil/internal/store/postgres"
)

func getEventStore(t *testing.T) *postgres.NotificationEventRepository {
	t.Helper()
	connectionString := os.Getenv("TEST_DB_CONNECTIONSTRING")
	if connectionString == "" {
		t.Skip("Skipping test since DB info not available")
		return nil
	}
	err := postgres.Migrate(connectionString)
	assert.Nil(t, err)
	dbc := postgres.NewStore(connectionString)
	return postgres.NewNotificationEventRepository(dbc)
}

func TestEvent(t *testing.T) {
	tearDown(t)
	db := getEventStore(t)
	ctx := context.Background()
	event := changedetector.NotificationEvent{
		ID:          "abc123",
		Type:        "SCHEMA_CHANGE_EVENT",
		EventTime:   time.Date(2024, 3, 15, 14, 18, 0, 0, time.UTC),
		NamespaceID: "gojek",
		SchemaID:    1,
		VersionID:   "version_id",
		Success:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	t.Run("NotificationEvents", func(t *testing.T) {
		t.Run("create: should create NotificationEvents", func(t *testing.T) {
			eventRes, err := db.Create(ctx, event)
			assert.Nil(t, err)
			assertEvent(t, event, eventRes)
		})
		t.Run("get: should get the by namespace, schema, version and success", func(t *testing.T) {
			eventRes, err := db.GetByNameSpaceSchemaVersionAndSuccess(ctx, "gojek", 1, "version_id", true)
			assert.Nil(t, err)
			assertEvent(t, event, eventRes)
		})
		t.Run("get: should update the event", func(t *testing.T) {
			eventRes, err := db.Update(ctx, "abc123", true)
			assert.Nil(t, err)
			assertEvent(t, event, eventRes)
		})
	})
}

func assertEvent(t *testing.T, expected, actual changedetector.NotificationEvent) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.EventTime, actual.EventTime)
	assert.Equal(t, expected.NamespaceID, actual.NamespaceID)
	assert.Equal(t, expected.SchemaID, actual.SchemaID)
	assert.Equal(t, expected.VersionID, actual.VersionID)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())
}
