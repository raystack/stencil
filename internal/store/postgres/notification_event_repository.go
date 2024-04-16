package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"

	"github.com/goto/stencil/core/changedetector"
)

const notificationEventsGetByNamespaceSchemaVersionSuccess = `
SELECT * from notification_events where namespace_id=$1 and schema_id=$2 and version_id=$3 and success=$4
`

const notificationEventsInsertQuery = `
INSERT INTO notification_events (id, type, event_time,namespace_id, schema_id, version_id, success,created_at,updated_at)
   VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *
`
const notificationEventsUpdateQuery = `
UPDATE notification_events SET updated_at=now(), success=$2
WHERE id = $1
RETURNING *
`

type NotificationEventRepository struct {
	db *DB
}

func NewNotificationEventRepository(dbc *DB) *NotificationEventRepository {
	return &NotificationEventRepository{
		db: dbc,
	}
}

func (r *NotificationEventRepository) Create(ctx context.Context, event changedetector.NotificationEvent) (changedetector.NotificationEvent, error) {
	newEvent := changedetector.NotificationEvent{}
	err := pgxscan.Get(ctx, r.db, &newEvent, notificationEventsInsertQuery, event.ID, event.Type, event.EventTime, event.NamespaceID, event.SchemaID,
		event.VersionID, event.Success)
	return newEvent, wrapError(err, event.NamespaceID, event.SchemaID, event.VersionID)
}

func (r *NotificationEventRepository) GetByNameSpaceSchemaVersionAndSuccess(ctx context.Context, namespace string, schemaID int32, versionID string, success bool) (changedetector.NotificationEvent, error) {
	newEvent := changedetector.NotificationEvent{}
	err := pgxscan.Get(ctx, r.db, &newEvent, notificationEventsGetByNamespaceSchemaVersionSuccess, namespace, schemaID, versionID, success)
	return newEvent, wrapError(err, namespace, schemaID, versionID, success)
}

func (r *NotificationEventRepository) Update(ctx context.Context, id string, success bool) (changedetector.NotificationEvent, error) {
	updatedEvent := changedetector.NotificationEvent{}
	err := pgxscan.Get(ctx, r.db, &updatedEvent, notificationEventsUpdateQuery, id, success)
	return updatedEvent, wrapError(err, id)
}
