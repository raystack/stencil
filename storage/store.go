package storage

import (
	"context"

	"github.com/odpf/stencil/models"
)

// Store is the interface that all database objects must implement.
type Store interface {
	// ListSnapshots returns a list of snapshots that match the query fields.
	ListSnapshots(ctx context.Context, query *models.Snapshot) ([]*models.Snapshot, error)

	// UpdateSnapshotLatestVersion updates the latest version of the snapshot.
	UpdateSnapshotLatestVersion(ctx context.Context, snapshot *models.Snapshot) error

	// GetSnapshotByFields returns the snapshot with the given fields.
	GetSnapshotByFields(ctx context.Context, namespace, name, version string, latest bool) (*models.Snapshot, error)

	// GetSnapshotByID returns the snapshot with the given id.
	GetSnapshotByID(ctx context.Context, id int64) (*models.Snapshot, error)

	// ExistsSnapshot returns true if the snapshot exists.
	ExistsSnapshot(ctx context.Context, snapshot *models.Snapshot) bool

	// CreateSnapshot creates a new snapshot.
	CreateSnapshot(ctx context.Context, snapshot *models.Snapshot) error

	// DeleteSnapshot deletes the snapshot.
	DeleteSnapshot(ctx context.Context, snapshot *models.Snapshot) error

	// GetSchema returns the protobuf file with the given name.
	GetSchema(ctx context.Context, snapshot *models.Snapshot, names []string) ([][]byte, error)

	// PutSchema puts the protobuf file with the given name.
	PutSchema(ctx context.Context, snapshot *models.Snapshot, dbFiles []*models.ProtobufDBFile) error
}
