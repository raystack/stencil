package api

import (
	"context"

	"github.com/odpf/stencil/server/snapshot"
)

//StoreService Service Interface for storage and validation
type StoreService interface {
	Validate(context.Context, *snapshot.Snapshot, []byte, []string) error
	Insert(context.Context, *snapshot.Snapshot, []byte) error
	Get(context.Context, *snapshot.Snapshot, []string) ([]byte, error)
}

// MetadataService Service Interface for metadata store
type MetadataService interface {
	Exists(context.Context, *snapshot.Snapshot) bool
	ListNames(context.Context, string) ([]string, error)
	ListVersions(context.Context, string, string) ([]string, error)
	GetSnapshot(context.Context, string, string, string, bool) (*snapshot.Snapshot, error)
	UpdateLatestVersion(context.Context, *snapshot.Snapshot) error
}

//API holds all handlers
type API struct {
	Store    StoreService
	Metadata MetadataService
}
