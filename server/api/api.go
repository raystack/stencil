package api

import (
	"context"

	"github.com/odpf/stencil/server/api/v1/genproto"
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
	List(context.Context, *snapshot.Snapshot) ([]*snapshot.Snapshot, error)
	GetSnapshotByFields(context.Context, string, string, string, bool) (*snapshot.Snapshot, error)
	GetSnapshotByID(context.Context, int64) (*snapshot.Snapshot, error)
	UpdateLatestVersion(context.Context, *snapshot.Snapshot) error
}

//API holds all handlers
type API struct {
	genproto.UnimplementedStencilServiceServer
	Store    StoreService
	Metadata MetadataService
}
